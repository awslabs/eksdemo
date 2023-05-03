package vpa

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type VpaOptions struct {
	application.ApplicationOptions

	AdmissionControllerEnabled bool
}

func newOptions() (options *VpaOptions, flags cmd.Flags) {
	options = &VpaOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.7.2",
				Latest:        "0.13.0",
				PreviousChart: "1.4.0",
				Previous:      "0.11.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "kube-system",
		},
		AdmissionControllerEnabled: false,
	}
	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admission-controller",
				Description: "enable VPA Admission Controller",
			},
			Option: &options.AdmissionControllerEnabled,
		},
	}
	return
}
