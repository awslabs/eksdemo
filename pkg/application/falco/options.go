package falco

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type FalcoOptions struct {
	*application.ApplicationOptions
	EventGenerator bool
	Replicas       int
}

func addOptions(a *application.Application) *application.Application {
	options := &FalcoOptions{
		ApplicationOptions: &application.ApplicationOptions{
			Namespace:      "falco",
			ServiceAccount: "falco",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.19.4",
				Latest:        "0.32.0",
				PreviousChart: "1.18.6",
				Previous:      "0.31.1",
			},
		},
		EventGenerator: false,
		Replicas:       1,
	}
	a.Options = options

	a.Flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "event-generator",
				Description: "enable the event generator deployment",
			},
			Option: &options.EventGenerator,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "replica count for event generator deployment",
			},
			Option: &options.Replicas,
		},
	}
	return a
}
