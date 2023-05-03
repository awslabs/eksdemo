package addon

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type AddonOptions struct {
	resource.CommonOptions

	Version string
}

func NewOptions() (options *AddonOptions, flags cmd.Flags) {
	options = &AddonOptions{}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "version",
				Description: "addon version (lookup default with \"eksdemo get addon-versions\")",
				Shorthand:   "v",
			},
			Option: &options.Version,
		},
	}
	return
}
