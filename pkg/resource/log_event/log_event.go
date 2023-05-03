package log_event

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "log-event",
			Description: "CloudWatch Log Events",
			Aliases:     []string{"log-events", "logevents", "logs", "log"},
			Args:        []string{"LOG_STREAM_NAME"},
		},

		Getter: &Getter{},
	}

	res.Options, res.GetFlags = newOptions()

	return res
}
