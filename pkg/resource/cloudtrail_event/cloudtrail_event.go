package cloudtrail_event

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cloudtrail-event",
			Description: "CloudTrail Event History",
			Aliases:     []string{"cloudtrail-events", "ctevents", "ctevent"},
			Args:        []string{"ID"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
