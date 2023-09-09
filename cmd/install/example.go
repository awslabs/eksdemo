package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/example/eks_workshop"
	"github.com/awslabs/eksdemo/pkg/application/example/game_2048"
	"github.com/awslabs/eksdemo/pkg/application/example/ghost"
	"github.com/awslabs/eksdemo/pkg/application/example/kube_ops_view"
	"github.com/awslabs/eksdemo/pkg/application/example/podinfo"
	"github.com/awslabs/eksdemo/pkg/application/example/wordpress"
	"github.com/spf13/cobra"
)

var exampleApps []func() *application.Application

func NewInstallExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "example",
		Short:   "Example Applications",
		Aliases: []string{"ex"},
	}

	// Don't show flag errors for `install example` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range exampleApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "example",
		Short:   "Example Applications",
		Aliases: []string{"ex"},
	}

	// Don't show flag errors for `uninstall example` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range exampleApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	exampleApps = []func() *application.Application{
		eks_workshop.NewApp,
		game_2048.NewApp,
		ghost.New,
		kube_ops_view.NewApp,
		podinfo.NewApp,
		wordpress.NewApp,
	}
}
