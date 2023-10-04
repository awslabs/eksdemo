package karpenter

import (
	"strings"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"

	"github.com/spf13/cobra"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	AMIFamily            string
	DisableDrift         bool
	Replicas             int
	TTLSecondsAfterEmpty int
}

func newOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.31.0",
				Latest:        "v0.31.0",
				PreviousChart: "v0.29.2",
				Previous:      "v0.29.2",
			},
		},
		AMIFamily: "AL2",
		Replicas:  1,
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ami-family",
				Description: "provisioner ami family",
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
				Description: "disables the drift deprovisioner",
			},
			Option: &options.DisableDrift,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ttl-after-empty",
				Description: "provisioner ttl seconds after empty (disables consolidation)",
				Shorthand:   "T",
			},
			Option: &options.TTLSecondsAfterEmpty,
		},
	}
	return
}
