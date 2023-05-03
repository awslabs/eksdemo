package listener_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type ListenerRuleOptions struct {
	resource.CommonOptions

	ListenerId       string
	LoadBalancerName string
}

func newOptions() (options *ListenerRuleOptions, flags cmd.Flags) {
	options = &ListenerRuleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "listener-id",
				Description: "listener id",
				Shorthand:   "I",
				Required:    true,
			},
			Option: &options.ListenerId,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "load balancer name",
				Shorthand:   "L",
				Required:    true,
			},
			Option: &options.LoadBalancerName,
		},
	}

	return
}
