package get

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/acm_certificate"
	"github.com/awslabs/eksdemo/pkg/resource/addon"
	"github.com/awslabs/eksdemo/pkg/resource/alarm"
	"github.com/awslabs/eksdemo/pkg/resource/amg_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/amp_rule"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/application"
	"github.com/awslabs/eksdemo/pkg/resource/auto_scaling_group"
	"github.com/awslabs/eksdemo/pkg/resource/availability_zone"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
	"github.com/awslabs/eksdemo/pkg/resource/cloudtrail_event"
	"github.com/awslabs/eksdemo/pkg/resource/cloudtrail_trail"
	"github.com/awslabs/eksdemo/pkg/resource/cluster"
	"github.com/awslabs/eksdemo/pkg/resource/dns_record"
	"github.com/awslabs/eksdemo/pkg/resource/ec2_instance"
	"github.com/awslabs/eksdemo/pkg/resource/ecr_repository"
	"github.com/awslabs/eksdemo/pkg/resource/elastic_ip"
	"github.com/awslabs/eksdemo/pkg/resource/event_rule"
	"github.com/awslabs/eksdemo/pkg/resource/fargate_profile"
	"github.com/awslabs/eksdemo/pkg/resource/hosted_zone"
	"github.com/awslabs/eksdemo/pkg/resource/iam_oidc"
	"github.com/awslabs/eksdemo/pkg/resource/iam_policy"
	"github.com/awslabs/eksdemo/pkg/resource/iam_role"
	"github.com/awslabs/eksdemo/pkg/resource/internet_gateway"
	"github.com/awslabs/eksdemo/pkg/resource/kms_key"
	"github.com/awslabs/eksdemo/pkg/resource/listener"
	"github.com/awslabs/eksdemo/pkg/resource/listener_rule"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
	"github.com/awslabs/eksdemo/pkg/resource/log_event"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
	"github.com/awslabs/eksdemo/pkg/resource/log_stream"
	"github.com/awslabs/eksdemo/pkg/resource/metric"
	"github.com/awslabs/eksdemo/pkg/resource/nat_gateway"
	"github.com/awslabs/eksdemo/pkg/resource/network_acl"
	"github.com/awslabs/eksdemo/pkg/resource/network_acl_rule"
	"github.com/awslabs/eksdemo/pkg/resource/network_interface"
	"github.com/awslabs/eksdemo/pkg/resource/node"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/resource/organization"
	"github.com/awslabs/eksdemo/pkg/resource/prefix_list"
	"github.com/awslabs/eksdemo/pkg/resource/route_table"
	"github.com/awslabs/eksdemo/pkg/resource/s3_bucket"
	"github.com/awslabs/eksdemo/pkg/resource/security_group"
	"github.com/awslabs/eksdemo/pkg/resource/security_group_rule"
	"github.com/awslabs/eksdemo/pkg/resource/sqs_queue"
	"github.com/awslabs/eksdemo/pkg/resource/ssm_node"
	"github.com/awslabs/eksdemo/pkg/resource/ssm_session"
	"github.com/awslabs/eksdemo/pkg/resource/subnet"
	"github.com/awslabs/eksdemo/pkg/resource/target_group"
	"github.com/awslabs/eksdemo/pkg/resource/target_health"
	"github.com/awslabs/eksdemo/pkg/resource/volume"
	"github.com/awslabs/eksdemo/pkg/resource/vpc"
	"github.com/awslabs/eksdemo/pkg/resource/vpc_endpoint"
	"github.com/awslabs/eksdemo/pkg/resource/vpc_summary"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "View resource(s)",
	}

	// Don't show flag errors for GET without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(acm_certificate.NewResource().NewGetCmd())
	cmd.AddCommand(addon.NewResource().NewGetCmd())
	cmd.AddCommand(addon.NewVersionsResource().NewGetCmd())
	cmd.AddCommand(alarm.NewResource().NewGetCmd())
	cmd.AddCommand(amg_workspace.NewResource().NewGetCmd())
	cmd.AddCommand(amp_rule.NewResource().NewGetCmd())
	cmd.AddCommand(amp_workspace.NewResource().NewGetCmd())
	cmd.AddCommand(application.NewResource().NewGetCmd())
	cmd.AddCommand(auto_scaling_group.NewResource().NewGetCmd())
	cmd.AddCommand(availability_zone.NewResource().NewGetCmd())
	cmd.AddCommand(cloudformation_stack.NewResource().NewGetCmd())
	cmd.AddCommand(cloudtrail_event.NewResource().NewGetCmd())
	cmd.AddCommand(cloudtrail_trail.NewResource().NewGetCmd())
	cmd.AddCommand(cluster.NewResource().NewGetCmd())
	cmd.AddCommand(NewGetCognitoCmd())
	cmd.AddCommand(NewGetAliasCmds(cognitoResources, "cognito-")...)
	cmd.AddCommand(dns_record.NewResource().NewGetCmd())
	cmd.AddCommand(ec2_instance.NewResource().NewGetCmd())
	cmd.AddCommand(ecr_repository.NewResource().NewGetCmd())
	cmd.AddCommand(elastic_ip.NewResource().NewGetCmd())
	cmd.AddCommand(event_rule.NewResource().NewGetCmd())
	cmd.AddCommand(fargate_profile.NewResource().NewGetCmd())
	cmd.AddCommand(hosted_zone.NewResource().NewGetCmd())
	cmd.AddCommand(iam_oidc.NewResource().NewGetCmd())
	cmd.AddCommand(iam_policy.NewResource().NewGetCmd())
	cmd.AddCommand(iam_role.NewResource().NewGetCmd())
	cmd.AddCommand(internet_gateway.NewResource().NewGetCmd())
	cmd.AddCommand(kms_key.NewResource().NewGetCmd())
	cmd.AddCommand(listener.NewResource().NewGetCmd())
	cmd.AddCommand(listener_rule.NewResource().NewGetCmd())
	cmd.AddCommand(load_balancer.NewResource().NewGetCmd())
	cmd.AddCommand(log_event.NewResource().NewGetCmd())
	cmd.AddCommand(log_group.NewResource().NewGetCmd())
	cmd.AddCommand(log_stream.NewResource().NewGetCmd())
	cmd.AddCommand(NewGetLogsInsightsCmd())
	cmd.AddCommand(NewGetAliasCmds(logInsights, "logs-insights-")...)
	cmd.AddCommand(NewGetAliasCmds(logInsights, "li-")...)
	cmd.AddCommand(metric.NewResource().NewGetCmd())
	cmd.AddCommand(nat_gateway.NewResource().NewGetCmd())
	cmd.AddCommand(network_acl.NewResource().NewGetCmd())
	cmd.AddCommand(network_acl_rule.NewResource().NewGetCmd())
	cmd.AddCommand(network_interface.NewResource().NewGetCmd())
	cmd.AddCommand(node.NewResource().NewGetCmd())
	cmd.AddCommand(nodegroup.NewResource().NewGetCmd())
	cmd.AddCommand(organization.NewResource().NewGetCmd())
	cmd.AddCommand(prefix_list.NewResource().NewGetCmd())
	cmd.AddCommand(route_table.NewResource().NewGetCmd())
	cmd.AddCommand(s3_bucket.NewResource().NewGetCmd())
	cmd.AddCommand(security_group.NewResource().NewGetCmd())
	cmd.AddCommand(security_group_rule.NewResource().NewGetCmd())
	cmd.AddCommand(sqs_queue.NewResource().NewGetCmd())
	cmd.AddCommand(ssm_node.NewResource().NewGetCmd())
	cmd.AddCommand(ssm_session.NewResource().NewGetCmd())
	cmd.AddCommand(subnet.NewResource().NewGetCmd())
	cmd.AddCommand(target_group.NewResource().NewGetCmd())
	cmd.AddCommand(target_health.NewResource().NewGetCmd())
	cmd.AddCommand(volume.NewResource().NewGetCmd())
	cmd.AddCommand(vpc.NewResource().NewGetCmd())
	cmd.AddCommand(vpc_endpoint.NewResource().NewGetCmd())
	cmd.AddCommand(NewGetVpcLatticeCmd())
	cmd.AddCommand(NewGetAliasCmds(vpcLattice, "vpc-lattice-")...)
	cmd.AddCommand(NewGetAliasCmds(vpcLattice, "vpclattice-")...)
	cmd.AddCommand(NewGetAliasCmds(vpcLattice, "lattice-")...)
	cmd.AddCommand(NewGetAliasCmds(vpcLattice, "vpcl-")...)
	cmd.AddCommand(vpc_summary.NewResource().NewGetCmd())

	return cmd
}

// This creates alias commands for subcommands under GET
func NewGetAliasCmds(resourceList []func() *resource.Resource, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(resourceList))

	for _, res := range resourceList {
		r := res()
		r.Command.Name = prefix + r.Command.Name
		r.Command.Hidden = true
		for i, alias := range r.Command.Aliases {
			r.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, r.NewGetCmd())
	}

	return cmds
}
