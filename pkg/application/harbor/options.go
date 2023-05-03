package harbor

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/spf13/cobra"
)

type HarborOptions struct {
	application.ApplicationOptions

	AdminPassword string
	NotaryEnabled bool
	NotaryHost    string
}

func newOptions() (options *HarborOptions, flags cmd.Flags) {
	options = &HarborOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.9.3",
				Latest:        "v2.5.3",
				PreviousChart: "1.9.3",
				Previous:      "v2.5.3",
			},
			DisableServiceAccountFlag: true,
			ExposeIngressOnly:         true,
			Namespace:                 "harbor",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "harbor admin password",
				Shorthand:   "P",
				Required:    true,
			},
			Option: &options.AdminPassword,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "notary",
				Description: "enable notary",
			},
			Option: &options.NotaryEnabled,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "notary-host",
				Description: "hostname for notary (required when Notary is enabled)",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.NotaryEnabled && options.NotaryHost == "" {
						return fmt.Errorf("%q flag required when using %q flag", "notary-host", "notary")
					}
					if !options.NotaryEnabled && options.NotaryHost != "" {
						return fmt.Errorf("%q flag required when using %q flag", "notary", "notary-host")
					}

					return nil
				},
			},
			Option: &options.NotaryHost,
		},
	}
	return
}
