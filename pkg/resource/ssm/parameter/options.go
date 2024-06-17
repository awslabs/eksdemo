package parameter

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Options struct {
	resource.CommonOptions

	// Get
	Path string
}

func newOptions() (options *Options, getFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "path",
				Description: "get parameters by path, instead of name",
				Validate: func(_ *cobra.Command, args []string) error {
					if options.Path != "" && len(args) > 0 {
						return &cmd.ArgumentAndFlagCantBeUsedTogetherError{Arg: "PARAMETER_NAME", Flag: "path"}
					}
					return nil
				},
			},
			Option: &options.Path,
		},
	}

	return
}
