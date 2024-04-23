package session

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type SessionOptions struct {
	resource.CommonOptions
	DocumentName string

	// Create
	PortForward      int
	PortForwardLocal int

	// Get
	History bool
}

func newOptions() (options *SessionOptions, createFlags, getFlags cmd.Flags) {
	options = &SessionOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
		DocumentName: "SSM-SessionManagerRunShell",
	}

	createFlags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "local-port",
				Description: "local port number when creating a port forwarding session",
				Shorthand:   "L",
				Validate: func(_ *cobra.Command, _ []string) error {
					if options.PortForwardLocal != 0 && options.PortForward == 0 {
						return &cmd.FlagRequiresFlagError{Flag1: "local-port", Flag2: "port-forward"}
					}
					if options.PortForwardLocal == 0 {
						options.PortForwardLocal = options.PortForward
					}
					return nil
				},
			},
			Option: &options.PortForwardLocal,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "port-forward",
				Description: "open a port forwarding session for the given port",
				Shorthand:   "P",
				Validate: func(_ *cobra.Command, _ []string) error {
					if options.PortForward != 0 {
						options.DocumentName = "AWS-StartPortForwardingSession"
					}
					return nil
				},
			},
			Option: &options.PortForward,
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
