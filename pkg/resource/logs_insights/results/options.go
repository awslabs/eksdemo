package results

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type ResultsFieldOptions struct {
	resource.CommonOptions

	Field     string
	LogStream string
	ShowStats bool
}

func newOptions() (options *ResultsFieldOptions, getFlags cmd.Flags) {
	options = &ResultsFieldOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
		Field: "@message",
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "field",
				Description: "results field to output values",
				Shorthand:   "f",
			},
			Option: &options.Field,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "log-stream",
				Description: "filter the results by log stream",
				Shorthand:   "l",
			},
			Option: &options.LogStream,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "stats",
				Description: "display query statistics after results",
			},
			Option: &options.ShowStats,
		},
	}

	return
}
