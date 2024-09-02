package domain

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	options, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "domain",
			Description: "SageMaker Domain",
			Args:        []string{"DOMAIN_NAME"},
		},

		GetFlags: getFlags,

		Getter: &Getter{},

		Options: options,
	}
}
