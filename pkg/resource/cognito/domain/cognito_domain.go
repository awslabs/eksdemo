package domain

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	options, createFlags, deleteFlags, getFlags := NewOptions()
	res := NewWithOptions(options)
	res.CreateFlags = createFlags
	res.DeleteFlags = deleteFlags
	res.GetFlags = getFlags

	return res
}

func NewWithOptions(options *Options) *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "domain",
			Description: "Cognito User Pool Domain",
			CreateArgs:  []string{"NAME"},
			Args:        []string{"DOMAIN"},
		},

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
