package userprofile

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Options struct {
	resource.CommonOptions

	// Get
	DomainID string
}

func newOptions() (options *Options, getFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                "sagemaker-domain",
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "domain-id",
				Description: "id of the sagemaker domain",
				Shorthand:   "D",
				Validate: func(_ *cobra.Command, args []string) error {
					if len(args) > 0 && options.DomainID != "" {
						return &cmd.ArgumentAndFlagCantBeUsedTogetherError{Arg: "USER_PROFILE_NAME", Flag: "domain-id"}
					}
					return nil
				},
			},
			Option: &options.DomainID,
		},
	}

	return
}
