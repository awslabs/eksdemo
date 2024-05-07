package ami

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	getFlags, options := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ami",
			Description: "Amazon Machine Image (AMI)",
			Aliases:     []string{"amis"},
			Args:        []string{"AMI_ID"},
		},

		GetFlags: getFlags,

		Getter: &Getter{},

		Options: options,
	}
}
