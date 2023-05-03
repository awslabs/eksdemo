package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/istio/istio_base"
	"github.com/awslabs/eksdemo/pkg/application/istio/istiod"
	"github.com/spf13/cobra"
)

var istioApps []func() *application.Application

func NewInstallIstioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Istio Service Mesh",
	}

	// Don't show flag errors for `install istio` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range istioApps {
		cmd.AddCommand(i().NewInstallCmd())
	}

	return cmd
}

func NewUninstallIstioCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "istio",
		Short: "Istio Service Mesh",
	}

	// Don't show flag errors for `uninstall istio` without a subcommand
	cmd.DisableFlagParsing = true

	for _, i := range istioApps {
		cmd.AddCommand(i().NewUninstallCmd())
	}

	return cmd
}

func init() {
	istioApps = []func() *application.Application{
		// bookinfo.NewApp,
		istio_base.NewApp,
		// istio_egress.NewApp,
		// istio_ingress.NewApp,
		istiod.NewApp,
		// kiali.NewApp,
	}
}
