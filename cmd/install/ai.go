package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/neuron/deviceplugin"
	"github.com/spf13/cobra"
)

var aiApps []func() *application.Application

func NewInstallAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI/ML Applications for Kubernetes",
	}

	// Don't show flag errors for `install ai` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range aiApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallAICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI/ML Applications for Kubernetes",
	}

	// Don't show flag errors for `uninstall ai` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range aiApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	aiApps = []func() *application.Application{
		deviceplugin.NewApp,
	}
}
