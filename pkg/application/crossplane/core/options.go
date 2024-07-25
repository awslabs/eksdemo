package core

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions
	ProviderVersion string
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.16.0",
				Latest:        "v1.16.0",
				PreviousChart: "1.16.0",
				Previous:      "v1.16.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "crossplane",
		},
		ProviderVersion: "v1.9.0",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "provider-version",
				Description: "version of provider-family-aws",
			},
			Option: &options.ProviderVersion,
		},
	}
	return
}
