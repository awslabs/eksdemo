package network_acl

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "network-acl",
			Description: "Network ACL",
			Aliases:     []string{"nacl"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	return res
}
