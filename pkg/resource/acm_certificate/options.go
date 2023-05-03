package acm_certificate

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type CertificateOptions struct {
	resource.CommonOptions

	sans           []string
	skipValidation bool
}

func NewOptions() (options *CertificateOptions, flags cmd.Flags) {
	options = &CertificateOptions{
		CommonOptions: resource.CommonOptions{
			Name:                "acm-certificate",
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "sans",
				Description: "subject alternative names",
			},
			Option: &options.sans,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "skip-validation",
				Description: "don't create Route53 record(s) to validate the certificate",
			},
			Option: &options.skipValidation,
		},
	}

	return
}
