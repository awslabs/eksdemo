package log_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, deleteFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "log-group",
			Description: "CloudWatch Log Group",
			Aliases:     []string{"log-groups", "loggroup", "lg"},
			CreateArgs:  []string{"NAME"},
			Args:        []string{"NAME"},
		},

		DeleteFlags: deleteFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
