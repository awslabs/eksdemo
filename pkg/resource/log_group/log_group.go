package log_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "log-group",
			Description: "CloudWatch Log Group",
			Aliases:     []string{"log-groups", "loggroup", "lg"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.GetFlags, res.DeleteFlags = newOptions()

	return res
}
