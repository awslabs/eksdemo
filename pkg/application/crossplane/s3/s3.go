package s3

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/crossplane"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	options := crossplane.NewProviderOptions("s3")

	return &application.Application{
		Command: cmd.Command{
			Parent:      "crossplane",
			Name:        "s3-provider",
			Description: "Crossplane S3 Provider",
			Aliases:     []string{"s3"},
		},

		Dependencies: []*resource.Resource{
			crossplane.CheckCore(),
			crossplane.Irsa(options, []string{"AmazonS3FullAccess"}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "crossplane-s3-provider",
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
  name: aws-s3
spec:
  serviceAccountTemplate:
    metadata:
      annotations:
        {{ .IrsaAnnotation }}         
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws-s3
spec:
  package: xpkg.upbound.io/upbound/provider-aws-s3:{{ .Version }}
  runtimeConfigRef:
    name: aws-s3
`
