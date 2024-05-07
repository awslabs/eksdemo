package ami

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Options struct {
	resource.CommonOptions

	NameFilter string
	Owners     []string
}

func newOptions() (getFlags cmd.Flags, o *Options) {
	o = &Options{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
		NameFilter: "amazon-eks-*",
		Owners:     []string{"amazon"},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "name",
				Description: "filter by the name of the AMI",
			},
			Option: &o.NameFilter,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "owners",
				Description: "scopes the results to images with the specified owners",
			},
			Option: &o.Owners,
		},
	}
	return
}
