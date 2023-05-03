package opa_gatekeeper

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type GatekeeperOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *GatekeeperOptions, flags cmd.Flags) {
	options = &GatekeeperOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace: "gatekeeper-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "3.9.0",
				Latest:        "v3.9.0",
				PreviousChart: "3.9.0",
				Previous:      "v3.9.0",
			},
			DisableServiceAccountFlag: true,
		},
		Replicas: 1,
	}
	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}
