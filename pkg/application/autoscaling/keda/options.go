package keda

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type KedaOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *KedaOptions, flags cmd.Flags) {
	options = &KedaOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.7.2",
				Latest:        "2.7.1",
				PreviousChart: "2.7.2",
				Previous:      "2.7.1",
			},
			Namespace:      "keda",
			ServiceAccount: "keda-operator",
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
