package network_interface

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type NetworkInterfaceOptions struct {
	resource.CommonOptions
	InstanceId       string
	IpAddress        string
	LoadBalancerName string
	SecurityGroupId  string
}

func NewOptions() (options *NetworkInterfaceOptions, flags cmd.Flags) {
	options = &NetworkInterfaceOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance-id",
				Description: "filter by Instance Id",
				Shorthand:   "I",
			},
			Option: &options.InstanceId,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ip-address",
				Description: "filter by IPv4 Address",
				Shorthand:   "A",
			},
			Option: &options.IpAddress,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "filter by Load Balancer name",
				Shorthand:   "L",
			},
			Option: &options.LoadBalancerName,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "security-group-id",
				Description: "filter by Security Group Id",
				Shorthand:   "S",
			},
			Option: &options.SecurityGroupId,
		},
	}
	return
}
