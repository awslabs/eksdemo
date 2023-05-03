package network_acl_rule

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type NetworkAclRuleOptions struct {
	resource.CommonOptions

	NetworkAclId string
	Egress       bool
}

func NewOptions() (options *NetworkAclRuleOptions, getFlags cmd.Flags) {
	options = &NetworkAclRuleOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nacl-id",
				Description: "network-acl id to display rules",
				Required:    true,
				Shorthand:   "N",
			},
			Option: &options.NetworkAclId,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "egress",
				Description: "show egress rules (defaults to ingress rules)",
			},
			Option: &options.Egress,
		},
	}

	return
}
