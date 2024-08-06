package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/csi/secretsstore"
	"github.com/spf13/cobra"
)

var secrets []func() *application.Application

func NewInstallSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secrets",
		Short:   "Secrets Management Solutions for Kubernetes",
		Aliases: []string{"secret"},
	}

	// Don't show flag errors for `install secrets` without a subcommand
	cmd.DisableFlagParsing = true

	for _, s := range secrets {
		cmd.AddCommand(s().NewInstallCmd())
	}

	return cmd
}

func NewUninstallSecretsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secrets",
		Short:   "Secrets Management Solutions for Kubernetes",
		Aliases: []string{"secret"},
	}

	// Don't show flag errors for `uninstall secrets` without a subcommand
	cmd.DisableFlagParsing = true

	for _, s := range secrets {
		cmd.AddCommand(s().NewUninstallCmd())
	}

	return cmd
}

func init() {
	secrets = []func() *application.Application{
		secretsstore.NewApp,
	}
}
