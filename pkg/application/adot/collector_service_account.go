package adot

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

func collectorServiceAccount(o *AdotOperatorOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &irsa.IrsaOptions{
			CommonOptions: resource.CommonOptions{
				ClusterName:    o.ClusterName,
				Name:           "adot-collector-service-account",
				Namespace:      o.Namespace,
				ServiceAccount: o.CollectorServiceAccount,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: saTemplate,
			},
		},
	}
	return res
}

const saTemplate = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  namespace: {{ .Namespace }}
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
