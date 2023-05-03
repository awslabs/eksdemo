package aws_fluent_bit

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type FluentbitOptions struct {
	*application.ApplicationOptions
	ReadFromTail bool
}

func addOptions(app *application.Application) *application.Application {
	options := &FluentbitOptions{
		ApplicationOptions: &application.ApplicationOptions{
			Namespace:      "logging",
			ServiceAccount: "fluent-bit",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.24.0",
				Latest:        "2.31.5",
				PreviousChart: "0.20.2",
				Previous:      "2.26.0",
			},
		},
	}
	app.Options = options

	app.Flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "read-from-tail",
				Description: "read the content from the tail of the file, not head",
			},
			Option: &options.ReadFromTail,
		},
	}
	return app
}
