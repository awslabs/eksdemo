package security_group

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type SecurityGroupOptions struct {
	resource.CommonOptions
	LoadBalancerName   string
	NetworkInterfaceId string
}

func NewOptions() (options *SecurityGroupOptions, flags cmd.Flags) {
	options = &SecurityGroupOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "eni-id",
				Description: "filter by Elastic Network Interface Id",
				Shorthand:   "E",
			},
			Option: &options.NetworkInterfaceId,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "filter by Load Balancer name",
				Shorthand:   "L",
			},
			Option: &options.LoadBalancerName,
		},
	}
	return
}
