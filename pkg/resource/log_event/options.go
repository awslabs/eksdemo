package log_event

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type LogEventOptions struct {
	resource.CommonOptions

	LogGroupName string
	Timestamp    bool
}

func newOptions() (options *LogEventOptions, getFlags cmd.Flags) {
	options = &LogEventOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	clusterFlag := options.NewClusterFlag(resource.Get, false)
	clusterFlag.Validate = nil

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "group-name",
				Description: "log group name",
				Shorthand:   "g",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.LogGroupName == "" && options.ClusterName == "" {
						return fmt.Errorf("must include either %q or %q flag", "--group-name", "--cluster")
					}

					if options.LogGroupName == "" {
						options.LogGroupName = fmt.Sprintf("/aws/eks/%s/cluster", options.ClusterName)
					}

					return nil
				},
			},
			Option: &options.LogGroupName,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "timestamp",
				Description: "include timestamp (only valid with \"table\" output)",
				Shorthand:   "T",
			},
			Option: &options.Timestamp,
		},
		clusterFlag,
	}

	return
}
