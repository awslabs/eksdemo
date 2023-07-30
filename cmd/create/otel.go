package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/otel"
	"github.com/spf13/cobra"
)

var otelCollectors []func() *resource.Resource

func NewOtelCollectorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "otel-collector",
		Short:   "Open Telemetry (OTEL) Collectors",
		Aliases: []string{"otel"},
	}

	// Don't show flag errors for `create otel-collector` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range otelCollectors {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	otelCollectors = []func() *resource.Resource{
		otel.NewPrometheusAMPCollector,
		otel.NewPrometheusCloudWatchCollector,
		otel.NewSimplestCollector,
	}
}
