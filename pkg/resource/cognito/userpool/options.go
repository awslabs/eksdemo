package userpool

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Options struct {
	resource.CommonOptions
	UserPoolName string
}

func NewOptions() (options *Options, createFlags, deleteFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                   "cognito-user-pool",
			DeleteArgumentOptional: true,
			ClusterFlagDisabled:    true,
		},
	}

	createFlags = cmd.Flags{}

	deleteFlags = cmd.Flags{}

	return
}

func (o *Options) SetName(name string) {
	o.UserPoolName = name
}
