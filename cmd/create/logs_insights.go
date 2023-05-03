package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/query"
	"github.com/spf13/cobra"
)

var logInsights []func() *resource.Resource

func NewCreateLogsInsightsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logs-insights",
		Short:   "CloudWatch Logs Insights Query",
		Aliases: []string{"li"},
	}

	// Don't show flag errors for `create logs-insights` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range logInsights {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	logInsights = []func() *resource.Resource{
		logs_insights.NewApiServerQuery,
		logs_insights.NewAuditQuery,
		logs_insights.NewAudit401Query,
		logs_insights.NewAudit403Query,
		logs_insights.NewAuthenticatorQuery,
		logs_insights.NewCloudControllerManagerQuery,
		logs_insights.NewControllerManagerQuery,
		logs_insights.NewControlPlaneQuery,
		logs_insights.NewSchedulerQuery,
		query.NewCreateResource,
	}
}
