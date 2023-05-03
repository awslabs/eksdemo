package crossplane

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

type AwsProviderOptions struct {
	irsa.IrsaOptions
	Version *string
}

func awsProvider(o *CrossplaneOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &AwsProviderOptions{
			IrsaOptions: irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					ClusterName:    o.Common().ClusterName,
					Name:           "crossplane-aws-provider",
					Namespace:      o.Common().Namespace,
					ServiceAccount: o.Common().ServiceAccount,
				},
			},
			Version: &o.ProviderVersion,
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

const yamlTemplate = `---
apiVersion: pkg.crossplane.io/v1alpha1
kind: ControllerConfig
metadata:
  name: aws-config
  annotations:
    {{ .IrsaAnnotation }}
spec:
  podSecurityContext:
    fsGroup: 2000
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws
spec:
  package: crossplane/provider-aws:{{ .Version }}
  controllerConfigRef:
    name: aws-config
`
