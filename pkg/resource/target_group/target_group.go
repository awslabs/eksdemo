package target_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "target-group",
			Description: "Target Group",
			Aliases:     []string{"target-groups", "tg"},
			CreateArgs:  []string{"NAME"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}
	res.Options, res.CreateFlags, res.GetFlags = newOptions()

	return res
}
