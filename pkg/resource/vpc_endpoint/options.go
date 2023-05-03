package vpc_endpoint

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type VpcEndpointOptions struct {
	resource.CommonOptions

	VpcId string
}

func newOptions() (options *VpcEndpointOptions, flags cmd.Flags) {
	options = &VpcEndpointOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vpc-id",
				Description: "filter by VPC Id",
				Shorthand:   "V",
			},
			Option: &options.VpcId,
		},
	}
	return
}
