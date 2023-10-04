package kube_prometheus_stack

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type KubePrometheusStackOptions struct {
	*application.ApplicationOptions
	GrafanaAdminPassword string
}

func newOptions() (options *KubePrometheusStackOptions, flags cmd.Flags) {
	options = &KubePrometheusStackOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "51.2.0",
				Latest:        "v0.68.0",
				PreviousChart: "46.6.0",
				Previous:      "v0.65.1",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "monitoring",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "grafana-pass",
				Description: "grafana admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.GrafanaAdminPassword,
		},
	}
	return
}
