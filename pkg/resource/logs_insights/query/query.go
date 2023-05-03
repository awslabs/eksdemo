package query

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewCreateResource() *resource.Resource {
	options, createFlags := NewOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "query",
			Description: "Logs Insights Query",
		},

		CreateFlags: createFlags,
		Manager:     &Manager{},
		Options:     options,
	}
}

func NewGetResource() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "query",
			Description: "Logs Insights Query History",
			Aliases:     []string{"queries", "history", "hist"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},

		Options: &resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}
}
