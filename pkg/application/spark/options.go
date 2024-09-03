package spark

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.0.0-rc.0",
				Latest:        "2.0.0-rc.0",
				PreviousChart: "2.0.0-rc.0",
				Previous:      "2.0.0-rc.0",
			},
			Namespace:      "spark-operator",
			ServiceAccount: "spark-operator-controller",
		},
	}

	flags = cmd.Flags{}

	return
}
