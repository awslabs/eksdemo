package delete

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/acm_certificate"
	"github.com/awslabs/eksdemo/pkg/resource/addon"
	"github.com/awslabs/eksdemo/pkg/resource/amg_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
	"github.com/awslabs/eksdemo/pkg/resource/cluster"
	"github.com/awslabs/eksdemo/pkg/resource/dns_record"
	"github.com/awslabs/eksdemo/pkg/resource/ec2_instance"
	"github.com/awslabs/eksdemo/pkg/resource/fargate_profile"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/resource/organization"
	"github.com/awslabs/eksdemo/pkg/resource/security_group"
	"github.com/awslabs/eksdemo/pkg/resource/target_group"
	"github.com/awslabs/eksdemo/pkg/resource/volume"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resource(s)",
	}

	// Don't show flag errors for delete without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(acm_certificate.NewResource().NewDeleteCmd())
	cmd.AddCommand(addon.NewResource().NewDeleteCmd())
	cmd.AddCommand(amg_workspace.NewResource().NewDeleteCmd())
	cmd.AddCommand(amp_workspace.NewResource().NewDeleteCmd())
	cmd.AddCommand(cloudformation_stack.NewResource().NewDeleteCmd())
	cmd.AddCommand(cluster.NewResource().NewDeleteCmd())
	cmd.AddCommand(NewCognitoCmd())
	cmd.AddCommand(NewDeleteAliasCmds(cognitoResources, "cognito-")...)
	cmd.AddCommand(dns_record.NewResource().NewDeleteCmd())
	cmd.AddCommand(ec2_instance.NewResource().NewDeleteCmd())
	cmd.AddCommand(fargate_profile.NewResource().NewDeleteCmd())
	cmd.AddCommand(irsa.NewResource().NewDeleteCmd())
	cmd.AddCommand(load_balancer.NewResource().NewDeleteCmd())
	cmd.AddCommand(log_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(nodegroup.NewResource().NewDeleteCmd())
	cmd.AddCommand(organization.NewResource().NewDeleteCmd())
	cmd.AddCommand(security_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(target_group.NewResource().NewDeleteCmd())
	cmd.AddCommand(volume.NewResource().NewDeleteCmd())

	return cmd
}

// This creates alias commands for subcommands under DELETE
func NewDeleteAliasCmds(resList []func() *resource.Resource, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(resList))

	for _, res := range resList {
		r := res()
		r.Command.Name = prefix + r.Command.Name
		r.Command.Hidden = true
		for i, alias := range r.Command.Aliases {
			r.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, r.NewDeleteCmd())
	}

	return cmds
}
