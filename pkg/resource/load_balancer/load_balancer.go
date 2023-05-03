package load_balancer

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "load-balancer",
			Description: "Elastic Load Balancer",
			Aliases:     []string{"load-balancers", "elbs", "elb", "lb"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},

		Options: &resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	return res
}
