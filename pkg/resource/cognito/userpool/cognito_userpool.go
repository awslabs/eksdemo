package userpool

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	options, createFlags, deleteFlags := NewOptions()
	res := NewWithOptions(options)
	res.CreateFlags = createFlags
	res.DeleteFlags = deleteFlags

	return res
}

func NewWithOptions(options *Options) *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "user-pool",
			Description: "Cognito User Pool",
			Aliases:     []string{"user-pools", "userpools", "userpool", "up"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
