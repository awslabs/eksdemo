package session

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, createFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-session",
			Description: "SSM Session",
			Aliases:     []string{"session"},
			Args:        []string{"INSTANCE_ID"},
		},

		CreateFlags: createFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
