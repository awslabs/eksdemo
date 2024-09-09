package cd

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
				LatestChart:   "7.5.2",
				Latest:        "v2.12.3",
				PreviousChart: "7.5.2",
				Previous:      "v2.12.3",
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
