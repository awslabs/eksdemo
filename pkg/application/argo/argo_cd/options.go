package argo_cd

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type ArgoCdOptions struct {
	application.ApplicationOptions

	AdminPassword string
}

func newOptions() (options *ArgoCdOptions, flags cmd.Flags) {
	options = &ArgoCdOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "6.10.0",
				Latest:        "v2.11.4",
				PreviousChart: "5.37.0",
				Previous:      "v2.7.7",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "argocd",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "Argo CD admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.AdminPassword,
		},
	}

	return
}
