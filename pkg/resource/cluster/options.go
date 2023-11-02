package cluster

import (
	"fmt"
	"net"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/autoscaling/karpenter"
	"github.com/awslabs/eksdemo/pkg/application/aws_lb_controller"
	"github.com/awslabs/eksdemo/pkg/application/external_dns"
	"github.com/awslabs/eksdemo/pkg/application/storage/ebs_csi"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

type ClusterOptions struct {
	resource.CommonOptions
	*nodegroup.NodegroupOptions

	DisableNetworkPolicy bool
	Fargate              bool
	HostnameType         string
	IPv6                 bool
	Kubeconfig           string
	NoRoles              bool
	PrefixAssignment     bool
	Private              bool
	VpcCidr              string

	appsForIrsa  []*application.Application
	IrsaTemplate *template.TextTemplate
	IrsaRoles    []*resource.Resource
}

func addOptions(res *resource.Resource) *resource.Resource {
	ngOptions, ngFlags, _ := nodegroup.NewOptions()

	options := &ClusterOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
			KubernetesVersion:   "1.28",
		},

		HostnameType:     string(types.HostnameTypeResourceName),
		NodegroupOptions: ngOptions,
		NoRoles:          false,
		VpcCidr:          "192.168.0.0/16",

		appsForIrsa: []*application.Application{
			aws_lb_controller.NewApp(),
			ebs_csi.NewApp(),
			external_dns.New(),
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

	// To keep in sync with ekctl, using logic from from DefaultPath() in eksctl/pkg/utils/kubeconfig
	if env := os.Getenv(clientcmd.RecommendedConfigPathEnvVar); len(env) > 0 {
		options.Kubeconfig = env
	} else {
		options.Kubeconfig = clientcmd.RecommendedHomeFile
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "version",
				Description: "Kubernetes version",
				Shorthand:   "v",
			},
			Choices: []string{"1.28", "1.27", "1.26", "1.25", "1.24"},
			Option:  &options.KubernetesVersion,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "disable-network-policy",
				Description: "don't enable network policy for Amazon VPC CNI",
			},
			Option: &options.DisableNetworkPolicy,
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
				Validate: func(cmd *cobra.Command, args []string) error {
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
				Validate: func(cmd *cobra.Command, args []string) error {
					_, _, err := net.ParseCIDR(options.VpcCidr)
					if err != nil {
						return fmt.Errorf("failed parsing --vpc-cidr, %w", err)
					}
					return nil
				},
			},
			Option: &options.VpcCidr,
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
