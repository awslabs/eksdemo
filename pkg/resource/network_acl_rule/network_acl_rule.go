package network_acl_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "network-acl-rule",
			Description: "Network ACL Rule",
			Aliases:     []string{"nacl-rules", "nacl-rule", "naclrules", "naclrule"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = NewOptions()

	return res
}
