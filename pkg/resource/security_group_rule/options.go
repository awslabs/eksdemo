package security_group_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type SecurityGroupRuleOptions struct {
	resource.CommonOptions
	Egress             bool
	Ingress            bool
	LoadBalancerName   string
	NetworkInterfaceId string
	SecurityGroupId    string
}

func NewOptions() (options *SecurityGroupRuleOptions, flags cmd.Flags) {
	options = &SecurityGroupRuleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "egress",
				Description: "show only egress rules",
			},
			Option: &options.Egress,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "eni-id",
				Description: "filter by Elastic Network Interface Id",
				Shorthand:   "E",
			},
			Option: &options.NetworkInterfaceId,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ingress",
				Description: "show only ingress rules",
			},
			Option: &options.Ingress,
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
