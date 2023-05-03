package ingress

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type BasicAuthSecretOptions struct {
	resource.CommonOptions
	AdminPassword string
}

func BasicAuthSecret(appName, password string) *resource.Resource {
	res := &resource.Resource{
		Options: &BasicAuthSecretOptions{
			CommonOptions: resource.CommonOptions{
				Name: appName + "-basic-auth-secret",
			},
			AdminPassword: password,
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: secretTemplate,
			},
		},
	}
	return res
}

const secretTemplate = `---
apiVersion: v1
kind: Secret
metadata:
  name: basic-auth
  namespace: {{ .Namespace }}
type: Opaque
data:
  auth: {{ htpasswd "admin" .AdminPassword | b64enc }}
`
