package ec2_instance

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "ec2-instance",
			Description: "EC2 Instance",
			Aliases:     []string{"ec2-instances", "ec2", "instances", "instance"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
