package target_health

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "target-health",
			Description: "Target Health",
			Aliases:     []string{"targets", "target", "th"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}
	res.Options, res.GetFlags = newOptions()

	return res
}
