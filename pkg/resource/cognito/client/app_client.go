package client

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
			Name:        "app-client",
			Description: "Cognito User Pool App Client",
			Args:        []string{"NAME"},
			Aliases:     []string{"client"},
		},

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
