package ghost

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/spf13/cobra"
)

type Options struct {
	application.ApplicationOptions

	StorageClass  string
	GhostPassword string
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "19.5.5",
				Latest:        "5.62.0",
				PreviousChart: "19.5.5",
				Previous:      "5.62.0",
			},
			ExposeIngressOnly: true,
			Namespace:         "ghost",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "storage-class",
				Description: "StorageClass for WordPress and MariaDB Persistent Volumes",
			},
			Option: &options.StorageClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ghost-pass",
				Description: "Ghost admin password",
				Required:    true,
				Shorthand:   "P",
				Validate: func(_ *cobra.Command, _ []string) error {
					if len(options.GhostPassword) >= 10 {
						return nil
					}
					return fmt.Errorf("the admin password must be at least 10 characters long")
				},
			},
			Option: &options.GhostPassword,
		},
	}
	return
}
