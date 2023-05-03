package iam_policy

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type IamPolicyOptions struct {
	resource.CommonOptions

	NameSearch string
	Role       string
}

func NewOptions() (options *IamPolicyOptions, getFlags cmd.Flags) {
	options = &IamPolicyOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "role",
				Description: "filter by role name, includes policy document in json/yaml output",
				Shorthand:   "r",
			},
			Option: &options.Role,
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
