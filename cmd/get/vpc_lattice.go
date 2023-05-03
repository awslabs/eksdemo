package get

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/vpc_lattice/service"
	"github.com/awslabs/eksdemo/pkg/resource/vpc_lattice/service_network"
	"github.com/awslabs/eksdemo/pkg/resource/vpc_lattice/target_group"
	"github.com/spf13/cobra"
)

var vpcLattice []func() *resource.Resource

func NewGetVpcLatticeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "vpc-lattice",
		Short:   "VPC Lattice Resources",
		Aliases: []string{"vpclattice", "lattice", "vpcl"},
	}

	// Don't show flag errors for `get vpc-lattice` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range vpcLattice {
		cmd.AddCommand(r().NewGetCmd())
	}

	return cmd
}

func init() {
	vpcLattice = []func() *resource.Resource{
		service.NewResource,
		service_network.NewResource,
		target_group.NewResource,
	}
}
