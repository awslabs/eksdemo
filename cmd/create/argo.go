package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/argo"
	"github.com/spf13/cobra"
)

var argoResources []func() *resource.Resource

func NewArgoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "argo",
		Short: "Argo Resources",
	}

	// Don't show flag errors for `create argo` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range argoResources {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	argoResources = []func() *resource.Resource{
		argo.NewGuestbook,
		argo.NewHelloWorld,
	}
}
