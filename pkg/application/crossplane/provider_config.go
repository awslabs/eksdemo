package crossplane

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func defaultProviderConfig() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "default-aws-provider-config",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: providerConfigManifest,
			},
		},
	}
	return res
}

const providerConfigManifest = `---
apiVersion: aws.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: InjectedIdentity
`
