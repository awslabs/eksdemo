package cloudtrail_trail

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cloudtrail-trail",
			Description: "CloudTrail Trail",
			Aliases:     []string{"cloudtrail-trails", "trails", "trail"},
			Args:        []string{"NAME_OR_ARN"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	return res
}
