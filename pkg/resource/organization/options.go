package organization

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type OrganizationOptions struct {
	resource.CommonOptions
}

func NewOptions() (options *OrganizationOptions, flags cmd.Flags) {
	options = &OrganizationOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "organization",
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{}

	return
}
