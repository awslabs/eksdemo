package iam_role

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"

	"github.com/spf13/cobra"
)

type IamRoleOptions struct {
	resource.CommonOptions

	LastUsed   bool
	NameSearch string
}

func NewOptions() (options *IamRoleOptions, getFlags cmd.Flags) {
	options = &IamRoleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	clusterFlag := options.NewClusterFlag(resource.Get, false)
	clusterFlag.Description = "filter by IRSA roles for cluster"
	origValidate := clusterFlag.Validate
	clusterFlag.Validate = func(cmd *cobra.Command, args []string) error {
		if options.ClusterName != "" && len(args) > 0 {
			return fmt.Errorf("%q flag cannot be used with NAME argument", "--cluster")
		}
		return origValidate(cmd, args)
	}

	getFlags = cmd.Flags{
		clusterFlag,
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "last-used",
				Description: "show last used date",
				Shorthand:   "L",
			},
			Option: &options.LastUsed,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "search",
				Description: "case-insensitive name search",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.NameSearch != "" && len(args) > 0 {
						return fmt.Errorf("%q flag cannot be used with NAME argument", "search")
					}
					return nil
				},
			},
			Option: &options.NameSearch,
		},
	}

	return
}
