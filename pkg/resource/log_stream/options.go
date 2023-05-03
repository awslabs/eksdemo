package log_stream

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"

	"github.com/spf13/cobra"
)

type LogStreamOptions struct {
	resource.CommonOptions

	LogGroupName string
}

func newOptions() (options *LogStreamOptions, getFlags cmd.Flags) {
	options = &LogStreamOptions{
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
						options.LogGroupName = log_group.LogGroupNameForClusterName(options.ClusterName)
					}

					return nil
				},
			},
			Option: &options.LogGroupName,
		},
		clusterFlag,
	}

	return
}
