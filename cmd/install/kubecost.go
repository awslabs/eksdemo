package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/kubecost/kubecost_eks"
	"github.com/awslabs/eksdemo/pkg/application/kubecost/kubecost_eks_amp"
	"github.com/awslabs/eksdemo/pkg/application/kubecost/kubecost_vendor"
	"github.com/spf13/cobra"
)

var kubecostApps []func() *application.Application

func NewInstallKubecostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubecost",
		Short: "Visibility Into Kubernetes Spend",
	}

	// Don't show flag errors for `install kubecost` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range kubecostApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallKubecostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubecost",
		Short: "Visibility Into Kubernetes Spend",
	}

	// Don't show flag errors for `uninstall kubecost` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range kubecostApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	kubecostApps = []func() *application.Application{
		kubecost_vendor.NewApp,
		kubecost_eks.NewApp,
		kubecost_eks_amp.NewApp,
	}
}
