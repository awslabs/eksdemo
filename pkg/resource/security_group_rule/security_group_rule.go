package security_group_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "security-group-rule",
			Description: "Security Group Rule",
			Aliases:     []string{"security-group-rules", "sg-rules", "sgrules", "sgr"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
