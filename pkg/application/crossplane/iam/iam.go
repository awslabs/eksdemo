package iam

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/crossplane"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	options := crossplane.NewProviderOptions("iam")

	return &application.Application{
		Command: cmd.Command{
			Parent:      "crossplane",
			Name:        "iam-provider",
			Description: "Crossplane IAM Provider",
			Aliases:     []string{"iam"},
		},

		Dependencies: []*resource.Resource{
			crossplane.CheckCore(),
			crossplane.Irsa(options, []string{"IAMFullAccess"}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "crossplane-iam-provider",
			ResourceTemplate: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},

		Options: options,
	}
}

const yamlTemplate = `---
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: aws-iam
spec:
  serviceAccountTemplate:
    metadata:
      annotations:
        {{ .IrsaAnnotation }}         
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws-iam
spec:
  package: xpkg.upbound.io/upbound/provider-aws-iam:{{ .Version }}
  runtimeConfigRef:
    name: aws-iam
`
