package vault

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AppOptions struct {
	application.ApplicationOptions
	Replicas int
}

func newOptions() (options *AppOptions, flags cmd.Flags) {
	options = &AppOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.28.0",
				Latest:        "v1.16.1",
				PreviousChart: "0.27.0",
				Previous:      "v1.15.2",
			},
			Namespace: "vault",
		},
		Replicas: 1,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "Number of replicas (3 recommended for High Availability)",
			},
			Option: &options.Replicas,
		},
	}

	return
}
