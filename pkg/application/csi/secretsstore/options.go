package secretsstore

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions
	RotateEnabled        bool
	RotationPollInterval string
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			Namespace: "kube-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.4",
				Latest:        "v1.4.4",
				PreviousChart: "1.4.4",
				Previous:      "v1.4.4",
			},
		},
		RotationPollInterval: "120s",
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "rotate",
				Description: "enables the secret rotation feature gate",
			},
			Option: &options.RotateEnabled,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "poll-interval",
				Description: "secret rotation poll interval duration",
			},
			Option: &options.RotationPollInterval,
		},
	}

	return
}
