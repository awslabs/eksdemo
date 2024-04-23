package session

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-session",
			Description: "SSM Session",
			Aliases:     []string{"session"},
			Args:        []string{"INSTANCE_ID"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}
	res.Options, res.GetFlags = newOptions()

	return res
}
