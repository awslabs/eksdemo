package dns_record

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "dns-record",
			Description: "Route53 Resource Record Set",
			Aliases:     []string{"dns-records", "records", "record", "dns"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.CreateFlags, res.DeleteFlags, res.GetFlags = newOptions()

	return res
}
