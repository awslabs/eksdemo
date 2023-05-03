package target_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "target-group",
			Description: "VPC Lattice Target Group",
			Aliases:     []string{"target-groups", "tg"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
