package instance

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, createFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ec2-instance",
			Description: "EC2 Instance",
			Aliases:     []string{"ec2-instances", "ec2", "instances", "instance"},
			Args:        []string{"INSTANCE_ID"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: createFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
