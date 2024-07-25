package ec2

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/crossplane"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	options := crossplane.NewProviderOptions("ec2")

	return &application.Application{
		Command: cmd.Command{
			Parent:      "crossplane",
			Name:        "ec2-provider",
			Description: "Crossplane EC2 Provider",
			Aliases:     []string{"ec2"},
		},

		Dependencies: []*resource.Resource{
			crossplane.CheckCore(),
			crossplane.Irsa(options, []string{"AmazonEC2FullAccess"}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "crossplane-ec2-provider",
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
  name: aws-ec2
spec:
  serviceAccountTemplate:
    metadata:
      annotations:
        {{ .IrsaAnnotation }}         
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws-ec2
spec:
  package: xpkg.upbound.io/upbound/provider-aws-ec2:{{ .Version }}
  runtimeConfigRef:
    name: aws-ec2
`
