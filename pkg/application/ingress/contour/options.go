package contour

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type ContourOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *ContourOptions, flags cmd.Flags) {
	options = &ContourOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "9.0.3",
				Latest:        "1.22.0",
				PreviousChart: "9.0.3",
				Previous:      "1.22.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "projectcontour",
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
