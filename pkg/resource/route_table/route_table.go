package route_table

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "route-table",
			Description: "Route Table",
			Aliases:     []string{"route-tables", "routetables", "routetable", "rt"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
