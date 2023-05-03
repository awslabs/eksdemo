package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/policy/kyverno"
	"github.com/awslabs/eksdemo/pkg/application/policy/opa_gatekeeper"
	"github.com/spf13/cobra"
)

var policyApps []func() *application.Application

func NewInstallPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Kubernetes Policy Controllers",
	}

	// Don't show flag errors for `install policy` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range policyApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallPolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Kubernetes Policy Controllers",
	}

	// Don't show flag errors for `uninstall policy` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range policyApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	policyApps = []func() *application.Application{
		kyverno.NewApp,
		opa_gatekeeper.NewApp,
	}
}
