package acm_certificate

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "acm-certificate",
			Description: "ACM Cerificate",
			Aliases:     []string{"acm-certificates", "acm-certs", "acm-cert", "acm"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options, res.CreateFlags = NewOptions()

	return res
}
