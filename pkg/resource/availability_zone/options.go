package availability_zone

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type AvailabilityZoneOptions struct {
	resource.CommonOptions
	AllZones bool
}

func newOptions() (options *AvailabilityZoneOptions, flags cmd.Flags) {
	options = &AvailabilityZoneOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "include all local zones and wavelength zones",
				Shorthand:   "A",
			},
			Option: &options.AllZones,
		},
	}

	return
}
