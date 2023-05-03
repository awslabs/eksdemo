package amg_workspace

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type AmgOptions struct {
	resource.CommonOptions

	Auth          []string
	Id            string
	WorkspaceName string
}

func NewOptions() (options *AmgOptions, createFlags, deleteFlags cmd.Flags) {
	options = &AmgOptions{
		CommonOptions: resource.CommonOptions{
			Name:                   "amazon-managed-grafana",
			DeleteArgumentOptional: true,
			ClusterFlagDisabled:    true,
		},
	}

	createFlags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "auth",
				Description: "Authentication methods (aws_sso, saml)",
				Required:    true,
				Validate: func(cmd *cobra.Command, args []string) error {
					for i, flag := range options.Auth {
						options.Auth[i] = strings.ToUpper(flag)
					}
					return nil
				},
			},
			Choices: []string{"aws_sso", "saml"},
			Option:  &options.Auth,
		},
	}

	deleteFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "delete by ID instead",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.Id == "" && len(args) == 0 {
						return fmt.Errorf("must include either %q argument or %q flag", "NAME", "--id")
					}
					return nil
				},
			},
			Option: &options.Id,
		},
	}

	return
}

func (o *AmgOptions) SetName(name string) {
	o.WorkspaceName = name
}

func (o *AmgOptions) iamRoleName() string {
	return fmt.Sprintf("eksdemo.amg.%s", o.WorkspaceName)
}
