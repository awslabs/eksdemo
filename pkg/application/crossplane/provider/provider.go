package provider

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/crossplane/core"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type Provider struct {
	Name            string
	Aliases         []string
	ManagedPolicies []string
}

func NewEC2() *application.Application {
	return NewProvider(&Provider{
		Name:            "ec2",
		Aliases:         []string{"ec2"},
		ManagedPolicies: []string{"AmazonEC2FullAccess"},
	})
}

func NewIAM() *application.Application {
	return NewProvider(&Provider{
		Name:            "iam",
		Aliases:         []string{"iam"},
		ManagedPolicies: []string{"IAMFullAccess"},
	})
}

func NewS3() *application.Application {
	return NewProvider(&Provider{
		Name:            "s3",
		Aliases:         []string{"s3"},
		ManagedPolicies: []string{"AmazonS3FullAccess"},
	})
}

func NewProvider(provider *Provider) *application.Application {
	options := newOptions(provider.Name)

	return &application.Application{
		Command: cmd.Command{
			Parent:      "crossplane",
			Name:        fmt.Sprintf("%s-provider", provider.Name),
			Description: fmt.Sprintf("Crossplane %s Provider", strings.ToUpper(provider.Name)),
			Aliases:     provider.Aliases,
		},

		Dependencies: []*resource.Resource{
			core.Check(),
			Irsa(options, provider.ManagedPolicies),
		},

		Installer: &installer.ManifestInstaller{
			AppName: fmt.Sprintf("crossplane-%s-provider", provider.Name),
			ResourceTemplate: &template.TextTemplate{
				Template: providerYamlTemplate,
			},
		},

		Options: options,
	}
}

const providerYamlTemplate = `---
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: aws-{{ .ProviderName }}
spec:
  serviceAccountTemplate:
    metadata:
      annotations:
        {{ .IrsaAnnotation }}         
---
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-aws-{{ .ProviderName }}
spec:
  package: xpkg.upbound.io/upbound/provider-aws-{{ .ProviderName }}:{{ .Version }}
  runtimeConfigRef:
    name: aws-{{ .ProviderName }}
`
