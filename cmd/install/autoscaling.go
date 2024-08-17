package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	"github.com/awslabs/eksdemo/pkg/application/autoscaling/keda"
	"github.com/awslabs/eksdemo/pkg/application/autoscaling/vpa"
	"github.com/spf13/cobra"
)

var autoscalingApps []func() *application.Application

func NewInstallAutoscalingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "autoscaling",
		Short:   "Kubernetes Autoscaling Applications",
		Aliases: []string{"as"},
	}

	// Don't show flag errors for `install autoscaling` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range autoscalingApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallAutoscalingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "autoscaling",
		Short:   "Kubernetes Autoscaling Applications",
		Aliases: []string{"as"},
	}

	// Don't show flag errors for `uninstall autoscaling` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range autoscalingApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	autoscalingApps = []func() *application.Application{
		cluster_autoscaler.NewApp,
		keda.NewApp,
		vpa.NewApp,
	}
}
