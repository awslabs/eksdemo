package karpenter

import (
	"strings"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"

	"github.com/spf13/cobra"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	AMIFamily        string
	DisableDrift     bool
	EnableSpotToSpot bool
	ExpireAfter      string
	Replicas         int
}

func newOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.34.2",
				Latest:        "v0.34.2",
				PreviousChart: "v0.33.2",
				Previous:      "v0.33.2",
			},
		},
		AMIFamily:   "AL2",
		ExpireAfter: "720h",
		Replicas:    1,
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ami-family",
				Description: "node class AMI family",
				Shorthand:   "A",
				Validate: func(cmd *cobra.Command, args []string) error {
					if strings.EqualFold(options.AMIFamily, "Al2") {
						options.AMIFamily = "AL2"
						return nil
					}
					if strings.EqualFold(options.AMIFamily, "Bottlerocket") {
						options.AMIFamily = "Bottlerocket"
						return nil
					}
					if strings.EqualFold(options.AMIFamily, "Ubuntu") {
						options.AMIFamily = "Ubuntu"
						return nil
					}
					return nil
				},
			},
			Option:  &options.AMIFamily,
			Choices: []string{"AL2", "Bottlerocket", "Ubuntu"},
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "disable-drift",
				Description: "disables the drift feature",
			},
			Option: &options.DisableDrift,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enable-spottospot",
				Description: "enables the spot to spot consolidation feature",
			},
			Option: &options.EnableSpotToSpot,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "expire-after",
				Description: "duration the controller will wait before terminating a node",
			},
			Option: &options.ExpireAfter,
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
