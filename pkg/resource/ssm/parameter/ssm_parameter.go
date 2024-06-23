package parameter

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-parameter",
			Description: "SSM Parameter",
			Aliases:     []string{"ssm-parameters", "ssm-params", "ssm-param", "params", "param"},
			Args:        []string{"PARAMETER_NAME"},
		},

		GetFlags: getFlags,

		Getter: &Getter{},

		Options: options,
	}
}
