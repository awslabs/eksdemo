package metric

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type CloudwatchMetricOptions struct {
	resource.CommonOptions

	Dimensions []string
	Namespace  string
	Values     bool

	cmd *cobra.Command
}

func newOptions() (options *CloudwatchMetricOptions, flags cmd.Flags) {
	options = &CloudwatchMetricOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "values",
				Description: "display dimension values (instead of metric counts)",
				Shorthand:   "V",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.Values && options.Namespace == "" {
						return fmt.Errorf("%q flag requires %q flag", "values", "namespace")
					}
					options.cmd = cmd
					return nil
				},
			},
			Option: &options.Values,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "dimensions",
				Description: "filter by dimension names",
				Shorthand:   "d",
				Validate: func(cmd *cobra.Command, args []string) error {
					if len(options.Dimensions) > 0 && options.Namespace == "" {
						return fmt.Errorf("%q flag requires %q flag", "dimension", "namespace")
					}
					return nil
				},
			},
			Option: &options.Dimensions,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "namespace",
				Description: "filter by metric namespace",
				Shorthand:   "n",
				Validate: func(cmd *cobra.Command, args []string) error {
					if len(args) > 0 && options.Namespace == "" {
						return fmt.Errorf("%q flag required with NAME argument", "namespace")
					}
					return nil
				},
			},
			Option: &options.Namespace,
		},
	}
	return
}
