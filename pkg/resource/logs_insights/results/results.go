package results

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "results",
			Description: "Logs Insights Query Results",
			Aliases:     []string{"result", "res"},
			Args:        []string{"QUERY_ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
			GetArgumentRequired: true,
		},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
