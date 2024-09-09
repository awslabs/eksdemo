package adot

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AdotOperatorOptions struct {
	application.ApplicationOptions
	CollectorServiceAccount string
}

func newOptions() (options *AdotOperatorOptions, flags cmd.Flags) {
	options = &AdotOperatorOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.63.2",
				Latest:        "0.102.0",
				PreviousChart: "0.31.0",
				Previous:      "0.78.0",
			},
			Namespace:      "adot-system",
			ServiceAccount: "adot-operator",
		},
		CollectorServiceAccount: "adot-collector",
	}

	return
}
