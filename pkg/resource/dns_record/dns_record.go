package dns_record

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, createFlags, deleteFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "dns-record",
			Description: "Route53 Resource Record Set",
			Aliases:     []string{"dns-records", "records", "record", "dns"},
			CreateArgs:  []string{"NAME"},
			Args:        []string{"NAME"},
		},

		CreateFlags: createFlags,
		DeleteFlags: deleteFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
