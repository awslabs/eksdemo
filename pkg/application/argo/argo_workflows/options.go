package argo_workflows

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type ArgoWorkflowOptions struct {
	application.ApplicationOptions

	AdminPassword string
	AuthMode      string
}

func newOptions() (options *ArgoWorkflowOptions, flags cmd.Flags) {
	options = &ArgoWorkflowOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.16.8",
				Latest:        "v3.3.8",
				PreviousChart: "0.16.8",
				Previous:      "v3.3.8",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "argo",
		},
		AuthMode: "server",
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
