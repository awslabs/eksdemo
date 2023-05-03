package node

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "node",
			Description: "Kubernetes Node",
			Aliases:     []string{"no", "nodes"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			Name: "kubernetes-node",
		},
	}

	return res
}
