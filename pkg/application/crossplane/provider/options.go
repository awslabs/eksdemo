package provider

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
)

type Options struct {
	application.ApplicationOptions
	ProviderName string
}

func newOptions(provider string) *Options {
	return &Options{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v1.9.0",
				Previous: "v1.8.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "crossplane",
			// Used only for role name in Crossplane IRSA
			ServiceAccount: fmt.Sprintf("provider-aws-%s", provider),
		},
		ProviderName: provider,
	}
}
