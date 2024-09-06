package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/ack"
	"github.com/awslabs/eksdemo/pkg/application/ack/apigatewayv2_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/ec2_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/ecr_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/eks_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/iam_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/prometheusservice_controller"
	"github.com/awslabs/eksdemo/pkg/application/ack/s3_controller"
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
		apigatewayv2_controller.NewApp,
		ec2_controller.NewApp,
		ecr_controller.NewApp,
		eks_controller.NewApp,
		iam_controller.NewApp,
		prometheusservice_controller.NewApp,
		ack.NewRDSController,
		s3_controller.NewApp,
	}
}
