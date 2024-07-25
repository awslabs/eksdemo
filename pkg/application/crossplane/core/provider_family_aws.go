package core

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type ProviderOptions struct {
	resource.CommonOptions
	Version *string
}

func providerFamilyAWS(o *Options) *resource.Resource {
	return &resource.Resource{
		Options: &ProviderOptions{
			CommonOptions: resource.CommonOptions{
				Name: "provider-family-aws",
			},
			Version: &o.ProviderVersion,
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: providerFamilyAWSManifest,
			},
		},
	}
}

const providerFamilyAWSManifest = `---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: upbound-provider-family-aws
spec:
  package: xpkg.upbound.io/upbound/provider-family-aws:{{ .Version }}
`
