package cluster

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/application"
	awslbc "github.com/awslabs/eksdemo/pkg/application/aws/lbc"
	"github.com/awslabs/eksdemo/pkg/application/externaldns"
	"github.com/awslabs/eksdemo/pkg/application/karpenter"
	"github.com/awslabs/eksdemo/pkg/application/storage/ebs_csi"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/kms/key"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

type ClusterOptions struct {
	resource.CommonOptions
	*nodegroup.NodegroupOptions

	Addons               []string
	DisableNetworkPolicy bool
	Fargate              bool
	HostnameType         string
	IPv6                 bool
	KMSKeyAlias          string
	KMSKeyArn            string
	Kubeconfig           string
	NoRoles              bool
	PrefixAssignment     bool
	Private              bool
	VpcCidr              string
	Zones                []string

	appsForIrsa  []*application.Application
	IrsaTemplate *template.TextTemplate
	IrsaRoles    []*resource.Resource
	Timeout      time.Duration
}

func addOptions(res *resource.Resource, resMgr *eksctl.ResourceManager) *resource.Resource {
	ngOptions, ngFlags, _ := nodegroup.NewOptions()

	options := &ClusterOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
			KubernetesVersion:   "1.32",
		},

		HostnameType:     string(types.HostnameTypeResourceName),
		Kubeconfig:       kubernetes.KubeconfigDefaultPath(),
		NodegroupOptions: ngOptions,
		NoRoles:          false,
		VpcCidr:          "192.168.0.0/16",

		appsForIrsa: []*application.Application{
			awslbc.NewApp(),
			ebs_csi.NewApp(),
			externaldns.New(),
			karpenter.NewApp(),
		},
		IrsaTemplate: &template.TextTemplate{
			Template: irsa.EksctlTemplate,
		},
	}

	ngOptions.CommonOptions = options.Common()
	ngOptions.DesiredCapacity = 2
	ngOptions.NodegroupName = "main"

	res.Options = options

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "version",
				Description: "Kubernetes version",
				Shorthand:   "v",
			},
			Choices: []string{"1.32", "1.31", "1.30", "1.29", "1.28", "1.27", "1.26", "1.25", "1.24"},
			Option:  &options.KubernetesVersion,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "disable-network-policy",
				Description: "don't enable network policy for Amazon VPC CNI",
			},
			Option: &options.DisableNetworkPolicy,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "encrypt-secrets",
				Description: "alias of KMS key to encrypt secrets",
				Validate: func(_ *cobra.Command, _ []string) error {
					if options.KMSKeyAlias == "" {
						return nil
					}
					key, err := key.NewGetter(aws.NewKMSClient()).GetByAlias(options.KMSKeyAlias)
					if err != nil {
						return err
					}
					options.KMSKeyArn = awssdk.ToString(key.Key.Arn)
					return nil
				},
			},
			Option: &options.KMSKeyAlias,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "fargate",
				Description: "create a Fargate profile",
			},
			Option: &options.Fargate,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "hostname-type",
				Description: "type of hostname to use for EC2 instances",
				Shorthand:   "H",
			},
			Choices: []string{string(types.HostnameTypeIpName), string(types.HostnameTypeResourceName)},
			Option:  &options.HostnameType,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ipv6",
				Description: "use IPv6 networking",
			},
			Option: &options.IPv6,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "kubeconfig",
				Description: "path to write kubeconfig",
				Validate: func(_ *cobra.Command, _ []string) error {
					// Set the KUBECONFIG environment variable to configure eksctl
					_ = os.Setenv("KUBECONFIG", options.Kubeconfig)
					return nil
				},
			},
			Option: &options.Kubeconfig,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-roles",
				Description: "don't create IAM roles",
			},
			Option: &options.NoRoles,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "prefix-assignment",
				Description: "configure VPC CNI for prefix assignment",
			},
			Option: &options.PrefixAssignment,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "private",
				Description: "private cluster (includes ECR, S3, and other VPC endpoints)",
			},
			Option: &options.Private,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vpc-cidr",
				Description: "CIDR to use for EKS Cluster VPC",
				Validate: func(_ *cobra.Command, _ []string) error {
					_, _, err := net.ParseCIDR(options.VpcCidr)
					if err != nil {
						return fmt.Errorf("failed parsing --vpc-cidr, %w", err)
					}
					return nil
				},
			},
			Option: &options.VpcCidr,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "zones",
				Description: "list of AZs to use. ie. us-east-1a,us-east-1b,us-east-1c",
			},
			Option: &options.Zones,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "addons",
				Description: "list of addons to use. examples: metrics-server,eks-pod-identity-agent,aws-ebs-csi-driver",
			},
			Option: &options.Addons,
		},
		&cmd.DurationFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "timeout",
				Description: "maximum waiting time for any long-running operation",
				Validate: func(_ *cobra.Command, _ []string) error {
					if options.Timeout.Seconds() > 0 {
						resMgr.CreateFlags = append(resMgr.CreateFlags,
							fmt.Sprintf("--timeout=%s", options.Timeout.String()))
					}
					return nil
				},
			},
			Option: &options.Timeout,
		},
	}

	res.CreateFlags = append(ngFlags, flags...)

	return res
}

func (o *ClusterOptions) PreCreate() error {
	o.Account = aws.AccountId()
	o.Partition = aws.Partition()
	o.Region = aws.Region()

	o.NodegroupOptions.IsClusterPrivate = o.Private
	o.NodegroupOptions.KubernetesVersion = o.KubernetesVersion

	// For apps we want to pre-create IRSA for, find the IRSA dependency
	for _, app := range o.appsForIrsa {
		for _, res := range app.Dependencies {
			if res.Name != "irsa" {
				continue
			}
			// Populate the IRSA Resource with Account, Cluster, Namespace, Partition, Region, ServiceAccount
			app.Common().Account = o.Account
			app.Common().ClusterName = o.ClusterName
			app.Common().Region = o.Region
			app.Common().Partition = o.Partition
			app.AssignCommonResourceOptions(res)
			res.SetName(app.Common().ServiceAccount)

			o.IrsaRoles = append(o.IrsaRoles, res)
		}
	}

	return o.NodegroupOptions.PreCreate()
}

func (o *ClusterOptions) PreDelete() error {
	o.Region = aws.Region()

	if err := karpenter.DeleteCustomResources(o.Common().KubeContext); err != nil {
		return err
	}

	cloudformationClient := aws.NewCloudformationClient()
	stacks, err := cloudformation_stack.NewGetter(cloudformationClient).GetStacksByCluster(o.ClusterName, "")
	if err != nil {
		return err
	}

	for _, stack := range stacks {
		stackName := awssdk.ToString(stack.StackName)
		if strings.HasPrefix(stackName, "eksdemo-") {
			fmt.Printf("Deleting Cloudformation stack %q\n", stackName)
			err := cloudformationClient.DeleteStack(stackName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *ClusterOptions) SetName(name string) {
	o.ClusterName = name
}
