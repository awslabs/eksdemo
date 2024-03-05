package accessentry

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "access-entry",
			Description: "EKS Access Entry",
			Aliases:     []string{"access-entries", "accessentry", "accessentries", "ae"},
			Args:        []string{"PRINCIPAL_ARN"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{},
	}
}
