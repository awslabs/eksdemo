package crossplane

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type CrossplaneOptions struct {
	application.ApplicationOptions
	ProviderVersion string
}

func newOptions() (options *CrossplaneOptions, flags cmd.Flags) {
	options = &CrossplaneOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.9.0",
				Latest:        "v1.9.0",
				PreviousChart: "1.9.0",
				Previous:      "v1.9.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "crossplane-system",
			// Used for role name in custom IRSA
			ServiceAccount: "provider-aws",
		},
		ProviderVersion: "v0.29.0",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "provider-version",
				Description: "version of provider-aws",
			},
			Option: &options.ProviderVersion,
		},
	}
	return
}
