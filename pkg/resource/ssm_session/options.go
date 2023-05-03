package ssm_session

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type SessionOptions struct {
	resource.CommonOptions

	History bool
	State   string
}

func newOptions() (options *SessionOptions, getFlags cmd.Flags) {
	options = &SessionOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "history",
				Description: "retrieve terminated sessions (instead of active)",
			},
			Option: &options.History,
		},
	}

	return
}
