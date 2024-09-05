package domain

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Options struct {
	resource.CommonOptions

	// Get, Delete
	DomainID string
}

func newOptions() (options *Options, deleteFlags, getFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled:    true,
			DeleteArgumentOptional: true,
		},
	}

	deleteFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "delete by id instead of name",
				Validate: func(_ *cobra.Command, args []string) error {
					if len(args) == 0 && options.DomainID == "" {
						return &cmd.MustIncludeEitherArgumentOrFlag{Arg: "DOMAIN_NAME", Flag: "--id"}
					}
					if options.DomainID != "" && len(args) > 0 {
						return &cmd.ArgumentAndFlagCantBeUsedTogetherError{Arg: "DOMAIN_NAME", Flag: "--id"}
					}
					return nil
				},
			},
			Option: &options.DomainID,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "get by id instead of name",
				Validate: func(_ *cobra.Command, args []string) error {
					if options.DomainID != "" && len(args) > 0 {
						return &cmd.ArgumentAndFlagCantBeUsedTogetherError{Arg: "DOMAIN_NAME", Flag: "--id"}
					}
					return nil
				},
			},
			Option: &options.DomainID,
		},
	}

	return
}
