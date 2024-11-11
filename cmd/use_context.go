package cmd

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var useContextCmd = &cobra.Command{
	Use:     "use-context CLUSTER",
	Short:   "Switch kubeconfig context to the EKS cluster specified",
	Aliases: []string{"context", "ctx", "uc"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterName := args[0]
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return aws.FormatErrorAsMessageOnly(err)
		}
		if cluster.Status != types.ClusterStatusActive {
			return fmt.Errorf("cluster %q is not active, cluster status is: %s", clusterName, cluster.Status)
		}

		configAccess := clientcmd.NewDefaultPathOptions()
		configFileName := configAccess.GetDefaultFilename()

		existingConfig, err := configAccess.GetStartingConfig()
		if err != nil {
			return fmt.Errorf("failed to read existing kubeconfig %q: %w", configFileName, err)
		}

		fmt.Printf("WARNING: The %q command has been deprecated. Please use the %q command.\n", "use-context", "use-cluster")

		contextName := contextForCluster(existingConfig, cluster)
		if contextName == "" {
			return fmt.Errorf("cluster %q not found in Kubeconfig", clusterName)
		}

		// update context
		existingConfig.CurrentContext = contextName

		err = clientcmd.ModifyConfig(configAccess, *existingConfig, false)
		if err != nil {
			return fmt.Errorf("failed to update kubeconfig: %w", err)
		}

		fmt.Printf("Context switched to: %s\n", contextName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(useContextCmd)
}
