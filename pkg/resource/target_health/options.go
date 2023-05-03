package target_health

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type TargetHealthOptions struct {
	resource.CommonOptions

	TargetGroupName string
}

func newOptions() (options *TargetHealthOptions, flags cmd.Flags) {
	options = &TargetHealthOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-group",
				Description: "name of Target Group",
				Shorthand:   "T",
				Required:    true,
			},
			Option: &options.TargetGroupName,
		},
	}

	return
}
