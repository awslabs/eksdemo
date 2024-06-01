package parameter

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ssm-parameter",
			Description: "SSM Parameter",
			Aliases:     []string{"ssm-parameters", "ssm-params", "ssm-param", "params", "param"},
			Args:        []string{"PATH_OR_NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}
}
