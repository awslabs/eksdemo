package consul

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type ConsulOptions struct {
	application.ApplicationOptions
	Namespace  string
	EnableUI   bool
	Replicas   int
	Datacenter string
}

func newOptions() (options *ConsulOptions, flags cmd.Flags) {
	options = &ConsulOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.1",
				Latest:        "v1.18.1",
				PreviousChart: "1.4.0",
				Previous:      "v1.18.0",
			},
		},
		Datacenter: "dc1",
		Namespace:  "consul",
		Replicas:   1,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enableUI",
				Description: "Enable Consul UI",
			},
			Option: &options.EnableUI,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "1 or 3 replicas",
			},
			Option: &options.Replicas,
		},
	}

	return
}
