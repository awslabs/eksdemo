package cloudformation_stack

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	return NewResourceWithOptions(newOptions())
}

func NewResourceWithOptions(options *CloudFormationOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cloudformation-stack",
			Description: "CloudFormation Stack",
			Aliases:     []string{"cloudformation-stacks", "stacks", "stack"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
