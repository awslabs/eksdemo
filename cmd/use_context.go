package cmd

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var useContextCmd = &cobra.Command{
	Use:     "use-context CLUSTER",
	Short:   "set currect kubeconfig context",
	Aliases: []string{"context", "ctx", "uc"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterName := args[0]
		eksctlClusterName := eksctl.GetClusterName(clusterName)

		config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)

		raw, _ := config.RawConfig()

		for name, context := range raw.Contexts {
			if context.Cluster == eksctlClusterName {
				// update context
				raw.CurrentContext = name

				err := clientcmd.ModifyConfig(config.ConfigAccess(), raw, false)
				if err != nil {
					return fmt.Errorf("failed to update kubeconfig: %w", err)
				}

				fmt.Printf("Context switched to: %s\n", name)
				return nil
			}
		}

		return fmt.Errorf("context not found in kubeconfig for cluster: %s", clusterName)
	},
}

func init() {
	rootCmd.AddCommand(useContextCmd)
}
