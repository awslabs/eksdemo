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
				// https://github.com/kubernetes/ingress-nginx#supported-versions-table
				// v1.8.0 supports k8s 1.27, 1.26, 1.25, 1.24
				LatestChart: "4.7.0",
				Latest:      "v1.8.0",
				// v1.6.x supports k8s 1.23
				PreviousChart: "4.5.2",
				Previous:      "v1.6.4",
			},
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
