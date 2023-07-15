package userpool

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Options struct {
	resource.CommonOptions
	UserPoolID   string
	UserPoolName string
}

func NewOptions() (options *Options, createFlags, deleteFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                   "cognito-user-pool",
			DeleteArgumentOptional: true,
			ClusterFlagDisabled:    true,
		},
	}

	createFlags = cmd.Flags{}

	deleteFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "delete by id instead",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.UserPoolID == "" && len(args) == 0 {
						return fmt.Errorf("must include either %q argument or %q flag", "NAME", "--id")
					}
					return nil
				},
			},
			Option: &options.UserPoolID,
		},
	}
	return
}

func (o *Options) SetName(name string) {
	o.UserPoolName = name
}
