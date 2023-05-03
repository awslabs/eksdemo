package auto_scaling_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "auto-scaling-group",
			Description: "Auto Scaling Group",
			Aliases:     []string{"auto-scaling-groups", "asgs", "asg"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
