package alarm

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "alarm",
			Description: "CloudWatch Alarm",
			Aliases:     []string{"alarms"},
			Args:        []string{"NAME_PREFIX"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}
}
