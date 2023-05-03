package resource

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func (o *CommonOptions) NewClusterFlag(action Action, required bool) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "cluster",
			Description: fmt.Sprintf("cluster to %s resource", action),
			Shorthand:   "c",
			Required:    required,
			Validate: func(cmd *cobra.Command, args []string) error {
				if !required && o.ClusterName == "" {
					return nil
				}

				cluster, err := aws.NewEKSClient().DescribeCluster(o.ClusterName)
				if err != nil {
					return aws.FormatErrorAsMessageOnly(err)
				}

				kubeContext, err := kubernetes.KubeContextForCluster(cluster)
				if err != nil {
					return err
				}
				if kubeContext == "" {
					return fmt.Errorf("cluster %q not found in Kubeconfig", o.ClusterName)
				}

				o.Cluster = cluster
				o.KubeContext = kubeContext
				o.KubernetesVersion = awssdk.ToString(cluster.Version)

				o.Account = aws.AccountId()
				o.Partition = aws.Partition()
				o.Region = aws.Region()

				return nil
			},
		},
		Option: &o.ClusterName,
	}
	return flag
}

func (o *CommonOptions) NewDryRunFlag() *cmd.BoolFlag {
	flag := &cmd.BoolFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "dry-run",
			Description: "don't create, just print out all creation steps",
		},
		Option: &o.DryRun,
	}
	return flag
}

func (o *CommonOptions) NewNamespaceFlag(action Action) *cmd.StringFlag {
	flag := &cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "namespace",
			Description: fmt.Sprintf("namespace to %s resource", action),
			Shorthand:   "n",
		},
		Option: &o.Namespace,
	}
	return flag
}
