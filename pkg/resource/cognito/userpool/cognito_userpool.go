package userpool

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "user-pool",
			Description: "Cognito User Pool",
			Aliases:     []string{"user-pools", "userpools", "userpool", "up"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}
}
