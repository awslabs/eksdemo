package karpenter_dashboards

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type KarpenterDashboardOptions struct {
	*application.ApplicationOptions
	KarpenterNamespace string
}

func newOptions() (options *KarpenterDashboardOptions, flags cmd.Flags) {
	options = &KarpenterDashboardOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "monitoring",
		},
		KarpenterNamespace: "karpenter",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "karpenter-namespace",
				Description: "namespace karpenter is installed",
			},
			Option: &options.KarpenterNamespace,
		},
	}
	return
}
