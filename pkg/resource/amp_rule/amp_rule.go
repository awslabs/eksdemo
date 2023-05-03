package amp_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amp-rule",
			Description: "Amazon Managed Prometheus Rule Namespace",
			Aliases:     []string{"amp-rules", "amprules", "amprule"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
