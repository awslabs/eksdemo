package vpc_endpoint

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc-endpoint",
			Description: "VPC Endpoint",
			Aliases:     []string{"vpc-endpoints", "endpoints", "endpoint"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
