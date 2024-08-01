package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/linkerd/base"
	"github.com/awslabs/eksdemo/pkg/application/linkerd/controlplane"
	"github.com/spf13/cobra"
)

var linkerdApps []func() *application.Application

func NewInstallLinkerdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "linkerd",
		Short: "Linkerd Service Mesh",
	}

	// Don't show flag errors for `install linkerd` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range linkerdApps {
		cmd.AddCommand( i().NewInstallCmd() )
	}

	return cmd
}

func NewUninstallLinkerdCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "linkerd",
		Short: "Linkerd Service Mesh",
	}

	// Don't show flag errors for `uninstall linkerd` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range linkerdApps {
		cmd.AddCommand(i().NewUninstallCmd())
	}

	return cmd
}

func init() {
	linkerdApps = []func() *application.Application{
		base.NewApp,
		controlplane.NewApp,
	}
}
