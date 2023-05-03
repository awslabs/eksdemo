package service_linked_role

import (
	"github.com/awslabs/eksdemo/pkg/resource"
)

type ServiceLinkedRoleOptions struct {
	resource.CommonOptions

	RoleName    string
	ServiceName string
}
