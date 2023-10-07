package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/acm_certificate"
	"github.com/awslabs/eksdemo/pkg/resource/addon"
	"github.com/awslabs/eksdemo/pkg/resource/amg_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/cluster"
	"github.com/awslabs/eksdemo/pkg/resource/dns_record"
	"github.com/awslabs/eksdemo/pkg/resource/fargate_profile"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/resource/organization"
	"github.com/awslabs/eksdemo/pkg/resource/ssm_session"
	"github.com/awslabs/eksdemo/pkg/resource/target_group"
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
	}

	// Don't show flag errors for create without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(NewAckCmd())
	cmd.AddCommand(NewCreateAliasCmds(ack, "ack-")...)
	cmd.AddCommand(acm_certificate.NewResource().NewCreateCmd())
	cmd.AddCommand(addon.NewResource().NewCreateCmd())
	cmd.AddCommand(amg_workspace.NewResource().NewCreateCmd())
	cmd.AddCommand(amp_workspace.NewResource().NewCreateCmd())
	cmd.AddCommand(NewArgoCmd())
	cmd.AddCommand(NewCreateAliasCmds(argoResources, "argo-")...)
	cmd.AddCommand(cluster.NewResource().NewCreateCmd())
	cmd.AddCommand(NewCognitoCmd())
	cmd.AddCommand(NewCreateAliasCmds(cognitoResources, "cognito-")...)
	cmd.AddCommand(dns_record.NewResource().NewCreateCmd())
	cmd.AddCommand(fargate_profile.NewResource().NewCreateCmd())
	cmd.AddCommand(NewKyvernoCmd())
	cmd.AddCommand(NewCreateAliasCmds(kyvernoPolicies, "kyverno-")...)
	cmd.AddCommand(log_group.NewResource().NewCreateCmd())
	cmd.AddCommand(NewCreateLogsInsightsCmd())
	cmd.AddCommand(NewCreateAliasCmds(logInsights, "logs-insights-")...)
	cmd.AddCommand(NewCreateAliasCmds(logInsights, "li-")...)
	cmd.AddCommand(nodegroup.NewResource().NewCreateCmd())
	cmd.AddCommand(nodegroup.NewSpotResource().NewCreateCmd())
	cmd.AddCommand(organization.NewResource().NewCreateCmd())
	cmd.AddCommand(NewOtelCollectorCmd())
	cmd.AddCommand(NewCreateAliasCmds(otelCollectors, "otel-collector-")...)
	cmd.AddCommand(NewCreateAliasCmds(otelCollectors, "otel-")...)
	cmd.AddCommand(ssm_session.NewResource().NewCreateCmd())
	cmd.AddCommand(target_group.NewResource().NewCreateCmd())

	return cmd
}

// This creates alias commands for subcommands under CREATE
func NewCreateAliasCmds(resList []func() *resource.Resource, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(resList))

	for _, res := range resList {
		r := res()
		r.Command.Name = prefix + r.Command.Name
		r.Command.Hidden = true
		for i, alias := range r.Command.Aliases {
			r.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, r.NewCreateCmd())
	}

	return cmds
}
