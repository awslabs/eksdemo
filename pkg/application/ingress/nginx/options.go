package nginx

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type NginxOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *NginxOptions, flags cmd.Flags) {
	options = &NginxOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "ingress-nginx",
			ServiceAccount: "ingress-nginx",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.4.2",
				Latest:        "v1.5.1",
				PreviousChart: "4.3.0",
				Previous:      "v1.4.0",
			},
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
