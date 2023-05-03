package inflate

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type InflateOptions struct {
	application.ApplicationOptions

	OnDemand bool
	Replicas int
	Spread   bool
}

func NewOptions() (options *InflateOptions, flags cmd.Flags) {
	options = &InflateOptions{
		ApplicationOptions: application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "inflate",
		},
		Replicas: 0,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "on-demand",
				Description: "request on-demand instances using karpenter node selector",
			},
			Option: &options.OnDemand,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the deployment (default 0)",
			},
			Option: &options.Replicas,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "spread",
				Description: "use topology spread constraints to spread across zones",
			},
			Option: &options.Spread,
		},
	}
	return
}
