package amp_workspace

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type AmpWorkspaceOptions struct {
	resource.CommonOptions

	Alias string
}

func NewOptions() (options *AmpWorkspaceOptions, flags cmd.Flags) {
	options = &AmpWorkspaceOptions{
		CommonOptions: resource.CommonOptions{
			Name:                   "amazon-managed-prometheus",
			DeleteArgumentOptional: true,
			ClusterFlagDisabled:    true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "delete by ID instead",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.Id == "" && len(args) == 0 {
						return fmt.Errorf("must include either %q argument or %q flag", "ALIAS", "--id")
					}
					return nil
				},
			},
			Option: &options.Id,
		},
	}

	return
}

func (o *AmpWorkspaceOptions) SetName(name string) {
	o.Alias = name
}
