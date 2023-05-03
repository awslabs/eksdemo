package get

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/query"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/results"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/stats"
	"github.com/spf13/cobra"
)

var logInsights []func() *resource.Resource

func NewGetLogsInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs-insights",
		Short:   "CloudWatch Logs Insights Query",
		Aliases: []string{"li"},
	}

	// Don't show flag errors for `get logs-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range logInsights {
		cmd.AddCommand(r().NewGetCmd())
	}

	return cmd
}

func init() {
	logInsights = []func() *resource.Resource{
		query.NewGetResource,
		results.NewResource,
		stats.NewResource,
	}
}
