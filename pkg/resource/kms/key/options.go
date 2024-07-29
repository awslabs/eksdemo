package key

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Options struct {
	resource.CommonOptions
}

func newOptions() (options *Options, getFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	getFlags = cmd.Flags{}

	return
}
