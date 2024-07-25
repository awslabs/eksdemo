package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/crossplane/s3"
	"github.com/awslabs/eksdemo/pkg/resource/crossplane/vpc"
	"github.com/spf13/cobra"
)

var crossplane []func() *resource.Resource

func NewCrossplaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crossplane",
		Short:   "The Cloud Native Control Plane",
		Aliases: []string{"cp"},
	}

	// Don't show flag errors for `create crossplane` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range crossplane {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	crossplane = []func() *resource.Resource{
		s3.NewResource,
		vpc.NewResource,
	}
}
