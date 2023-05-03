package application

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "application",
			Description: "Installed Applications",
			Aliases:     []string{"applications", "apps", "app"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{},
	}

	return res
}
