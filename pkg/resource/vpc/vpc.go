package vpc

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc",
			Description: "Virtual Private Cloud",
			Aliases:     []string{"vpcs"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	return res
}
