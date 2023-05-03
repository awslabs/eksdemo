package cloudformation_stack

import (
	"github.com/awslabs/eksdemo/pkg/resource"
)

type CloudFormationOptions struct {
	resource.CommonOptions
}

func newOptions() (options *CloudFormationOptions) {
	options = &CloudFormationOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "cloudformation",
			ClusterFlagDisabled: true,
			ClusterFlagOptional: true,
		},
	}

	return
}
