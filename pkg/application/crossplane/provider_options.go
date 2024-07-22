package crossplane

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
)

type ProviderOptions struct {
	application.ApplicationOptions
	ProviderName string
}

func NewProviderOptions(provider string) *ProviderOptions {
	providerName := fmt.Sprintf("provider-aws-%s", provider)

	return &ProviderOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				Latest:   "v1.9.0",
				Previous: "v1.8.0",
			},
			DisableServiceAccountFlag: true,
			Namespace:                 "crossplane",
			// Used only for role name in Crossplane IRSA
			ServiceAccount: providerName,
		},
		ProviderName: providerName,
	}
}
