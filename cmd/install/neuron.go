package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/neuron/deviceplugin"
	"github.com/spf13/cobra"
)

var neuron []func() *application.Application

func NewInstallNeuronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neuron",
		Short: "AWS Neuron for Inferentia and Trainium Support",
	}

	// Don't show flag errors for `install ai` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range neuron {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallNeuronCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neuron",
		Short: "AWS Neuron for Inferentia and Trainium Support",
	}

	// Don't show flag errors for `uninstall ai` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range neuron {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	neuron = []func() *application.Application{
		deviceplugin.NewApp,
	}
}
