package stats

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "stats",
			Description: "Logs Insights Query Statistics",
			Aliases:     []string{"statistics", "statistic", "stat"},
			Args:        []string{"QUERY_ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
			GetArgumentRequired: true,
		},
	}

	return res
}
