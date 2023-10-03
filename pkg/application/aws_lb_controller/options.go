package aws_lb_controller

import (
	"strings"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/spf13/cobra"
)

type AWSLBControllerOptions struct {
	application.ApplicationOptions

	DefaultIngressClass bool
	DefaultTargetType   string
	DisableWebhook      bool
	Replicas            int
}

func newOptions() (options *AWSLBControllerOptions, flags cmd.Flags) {
	options = &AWSLBControllerOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "awslb",
			ServiceAccount: "aws-load-balancer-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.6.1",
				Latest:        "v2.6.1",
				PreviousChart: "1.6.0",
				Previous:      "v2.6.0",
			},
		},
		DefaultTargetType: "instance",
		Replicas:          1,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "default-ingress-class",
				Description: "set the alb IngressClass as the default for the cluster",
			},
			Option: &options.DefaultIngressClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "default-target-type",
				Description: "set the default target type for target groups",
				Validate: func(cmd *cobra.Command, args []string) error {
					options.DefaultTargetType = strings.ToLower(options.DefaultTargetType)
					return nil
				},
			},
			Option:  &options.DefaultTargetType,
			Choices: []string{"instance", "ip"},
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "disable-webhook",
				Description: "disable the webhook so the in-tree controller can provision CLBs",
			},
			Option: &options.DisableWebhook,
		},
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
