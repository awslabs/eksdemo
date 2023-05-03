package vpc_summary

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type VpcSummaryOptions struct {
	resource.CommonOptions
	ShowIds bool
}

func newOptions() (options *VpcSummaryOptions, flags cmd.Flags) {
	options = &VpcSummaryOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ids",
				Description: "show ids",
			},
			Option: &options.ShowIds,
		},
	}
	return
}
