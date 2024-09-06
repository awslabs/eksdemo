package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/ack"
	"github.com/spf13/cobra"
)

var ackControllers []func() *application.Application

func NewInstallAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `install ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range ackControllers {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `uninstall ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range ackControllers {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	ackControllers = []func() *application.Application{
		ack.NewAPIGatewayv2Controller,
		ack.NewEC2Controller,
		ack.NewECRController,
		ack.NewEFSController,
		ack.NewEKSController,
		ack.NewIAMController,
		ack.NewPrometheusServiceController,
		ack.NewRDSController,
		ack.NewS3Controller,
	}
}
