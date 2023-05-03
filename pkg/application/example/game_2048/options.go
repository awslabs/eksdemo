package game_2048

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Game2048Options struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *Game2048Options, flags cmd.Flags) {
	options = &Game2048Options{
		ApplicationOptions: application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "game-2048",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "latest",
				Previous: "latest",
			},
			DisableServiceAccountFlag: true,
		},
		Replicas: 1,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}
