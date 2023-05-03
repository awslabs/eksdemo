package service_linked_role

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResourceWithOptions(options *ServiceLinkedRoleOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "service-linked-role",
			Description: "Service Linked Role",
		},

		Manager: &Manager{},

		Options: options,
	}

	return res
}
