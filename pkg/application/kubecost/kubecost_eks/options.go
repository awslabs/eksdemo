package kubecost_eks

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type KubecostEksOptions struct {
	application.ApplicationOptions

	EnableNodeExporter bool
}

func newOptions() (options *KubecostEksOptions, flags cmd.Flags) {
	options = &KubecostEksOptions{
		ApplicationOptions: application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kubecost",
			ServiceAccount:               "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.103.3",
				Latest:        "1.103.3",
				PreviousChart: "1.100.0",
				Previous:      "1.100.0",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "node-exporter",
				Description: "install prometheus node exporter (not installed by default)",
			},
			Option: &options.EnableNodeExporter,
		},
	}
	return
}
