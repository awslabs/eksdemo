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
				LatestChart:   "4.9.14",
				Latest:        "v2.4.6",
				PreviousChart: "4.9.12",
				Previous:      "v2.4.4",
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
