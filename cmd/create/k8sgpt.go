package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/k8sgpt/bedrock/claudev2"
	"github.com/spf13/cobra"
)

var k8sgpt []func() *resource.Resource

func NewK8sGPTCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "k8sgpt",
		Short: "K8sGPT Resources",
	}

	// Don't show flag errors for `create k8sgpt` without a subcommand
	cmd.DisableFlagParsing = true

	for _, k := range k8sgpt {
		cmd.AddCommand(k().NewCreateCmd())
	}

	return cmd
}

func init() {
	k8sgpt = []func() *resource.Resource{
		claudev2.NewResource,
	}
}
