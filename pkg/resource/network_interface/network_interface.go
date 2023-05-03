package network_interface

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "network-interface",
			Description: "Elastic Network Interface",
			Aliases:     []string{"network-interfaces", "enis", "eni"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
