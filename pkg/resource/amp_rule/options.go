package amp_rule

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type AmpRuleOptions struct {
	resource.CommonOptions

	Alias       string
	WorkspaceId string
}

func newOptions() (options *AmpRuleOptions, flags cmd.Flags) {
	options = &AmpRuleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "alias",
				Description: "AMP workspace alias",
				Shorthand:   "a",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.Alias == "" && options.WorkspaceId == "" {
						return fmt.Errorf("either %q flag or %q flag is required", "alias", "workspace-id")
					}
					if options.Alias != "" && options.WorkspaceId != "" {
						return fmt.Errorf("%q flag and %q flag can not be used together", "alias", "workspace-id")
					}
					return nil
				},
			},
			Option: &options.Alias,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "workspace-id",
				Description: "AMP workspace id",
				Shorthand:   "W",
			},
			Option: &options.WorkspaceId,
		},
	}

	return
}
