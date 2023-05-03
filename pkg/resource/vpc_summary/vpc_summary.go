package vpc_summary

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, flags := newOptions()

	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc-summary",
			Description: "VPC Summary",
			Aliases:     []string{"vpcsummary", "vpcsum"},
			Args:        []string{"ID"},
		},

		GetFlags: flags,

		Getter: &Getter{},

		Options: options,
	}

	return res
}
