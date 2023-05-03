package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/container_insights/adot_collector"
	"github.com/awslabs/eksdemo/pkg/application/container_insights/cloudwatch_agent"
	"github.com/awslabs/eksdemo/pkg/application/container_insights/fluent_bit"
	"github.com/awslabs/eksdemo/pkg/application/container_insights/prometheus"
	"github.com/spf13/cobra"
)

var containerInsightsApps []func() *application.Application

func NewInstallContainerInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "container-insights",
		Short:   "CloudWatch Container Insights",
		Aliases: []string{"ci"},
	}

	// Don't show flag errors for `install container-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range containerInsightsApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallContainerInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "container-insights",
		Short:   "CloudWatch Container Insights",
		Aliases: []string{"ci"},
	}

	// Don't show flag errors for `uninstall container-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range containerInsightsApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	containerInsightsApps = []func() *application.Application{
		adot_collector.NewApp,
		cloudwatch_agent.NewApp,
		fluent_bit.NewApp,
		prometheus.NewApp,
	}
}
