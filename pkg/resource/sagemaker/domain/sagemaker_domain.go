package domain

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	options, deleteFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "domain",
			Description: "SageMaker Domain",
			Aliases:     []string{"do"},
			Args:        []string{"DOMAIN_NAME"},
		},

		DeleteFlags: deleteFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
