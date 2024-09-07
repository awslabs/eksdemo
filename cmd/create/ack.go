package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/ack/amp"
	"github.com/awslabs/eksdemo/pkg/resource/ack/ec2"
	"github.com/awslabs/eksdemo/pkg/resource/ack/ecr"
	"github.com/awslabs/eksdemo/pkg/resource/ack/efs"
	"github.com/awslabs/eksdemo/pkg/resource/ack/eks"
	"github.com/awslabs/eksdemo/pkg/resource/ack/iam"
	"github.com/awslabs/eksdemo/pkg/resource/ack/rds"
	"github.com/awslabs/eksdemo/pkg/resource/ack/s3"
	"github.com/spf13/cobra"
)

var ack []func() *resource.Resource

func NewAckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ack",
		Short: "AWS Controllers for Kubernetes (ACK)",
	}

	// Don't show flag errors for `create ack` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range ack {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	ack = []func() *resource.Resource{
		amp.NewLoggingConfigurationResource,
		amp.NewWorkspaceResource,
		ec2.NewSecurityGroupResource,
		ec2.NewSubnetResource,
		ec2.NewVpcResource,
		efs.NewFileSystemResource,
		ecr.NewResource,
		iam.NewRoleResource,
		eks.NewFargateProfileResource,
		rds.NewDatabaseInstanceResource,
		s3.NewResource,
	}
}
