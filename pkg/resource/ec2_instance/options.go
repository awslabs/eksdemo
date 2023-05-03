package ec2_instance

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type EC2InstanceOptions struct {
	resource.CommonOptions

	HideTerminated bool
}

func newOptions() (o *EC2InstanceOptions, flags cmd.Flags) {
	o = &EC2InstanceOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "hide-terminated",
				Description: "show only pending, running, shutting-down, stopping, stopped instances",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.HideTerminated && len(args) > 0 {
						return fmt.Errorf("%q flag cannot be used with ID argument", "hide-terminated")
					}
					return nil
				},
			},
			Option: &o.HideTerminated,
		},
	}
	return
}
