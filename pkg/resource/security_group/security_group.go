package security_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "security-group",
			Description: "Security Group",
			Aliases:     []string{"security-groups", "sg"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
