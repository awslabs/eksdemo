package domain

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Options struct {
	resource.CommonOptions
	DomainName string
}

func NewOptions() (options *Options, createFlags, deleteFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                "cognito-domain",
			ClusterFlagDisabled: true,
			GetArgumentRequired: true,
		},
	}

	createFlags = cmd.Flags{}

	deleteFlags = cmd.Flags{}

	return
}

func (o *Options) SetName(name string) {
	o.DomainName = name
}
