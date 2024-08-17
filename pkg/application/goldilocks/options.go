package goldilocks

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions
	NoVPA bool
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "8.0.2",
				Latest:        "v4.13.0",
				PreviousChart: "8.0.2",
				Previous:      "v4.13.0",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "goldilocks",
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-vpa",
				Description: "don't install the VPA sub-chart",
			},
			Option: &options.NoVPA,
		},
	}

	return
}
