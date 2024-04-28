package consul

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AppOptions struct {
	application.ApplicationOptions
	EnableUI    bool
	EnableMesh  bool
	EnableAPIGW bool
	Replicas    int
	Datacenter  string
}

func newOptions() (options *AppOptions, flags cmd.Flags) {
	options = &AppOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.1",
				Latest:        "v1.18.1",
				PreviousChart: "1.4.0",
				Previous:      "v1.18.0",
			},
			Namespace: "consul",
		},
		Datacenter: "dc1",
		Replicas:   1,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enable-ui",
				Description: "Enable Consul UI",
			},
			Option: &options.EnableUI,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enable-mesh",
				Description: "Enable Consul Service Mesh",
			},
			Option: &options.EnableMesh,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enable-api-gw",
				Description: "Enable Consul API Gateway",
			},
			Option: &options.EnableAPIGW,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "datacenter",
				Description: "Specify Consul Datacenter",
			},
			Option: &options.Datacenter,
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
