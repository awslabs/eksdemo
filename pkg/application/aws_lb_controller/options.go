package aws_lb_controller

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AWSLBControllerOptions struct {
	application.ApplicationOptions

	Default bool
}

func newOptions() (options *AWSLBControllerOptions, flags cmd.Flags) {
	options = &AWSLBControllerOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "awslb",
			ServiceAccount: "aws-load-balancer-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.8",
				Latest:        "v2.4.7",
				PreviousChart: "1.4.7",
				Previous:      "v2.4.6",
			},
		},
		Default: false,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "default",
				Description: "set as the default IngressClass for the cluster",
			},
			Option: &options.Default,
		},
	}

	return
}
