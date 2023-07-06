package workflows

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions

	AuthMode string
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.32.1",
				Latest:        "v3.4.9",
				PreviousChart: "0.16.8",
				Previous:      "v3.3.8",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "argo",
		},
		AuthMode: "client",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "auth-mode",
				Description: "Argo Server authentication mode",
			},
			Option:  &options.AuthMode,
			Choices: []string{"client", "server", "sso"},
		},
	}

	return
}
