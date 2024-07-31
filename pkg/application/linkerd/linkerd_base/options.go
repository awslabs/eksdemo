package linkerdBase

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
				LatestChart:   "2024.7.3",
				PreviousChart: "2024.7.2",
			},
			Namespace: "linkerd",
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
