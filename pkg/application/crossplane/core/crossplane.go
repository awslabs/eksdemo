package core

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.crossplane.io
// GitHub:  https://github.com/crossplane/crossplane
// Helm:    https://github.com/crossplane/crossplane/tree/master/cluster/charts/crossplane
// Repo:    xpkg.upbound.io/crossplane/crossplane
// Version: Latest is Chart/App v1.16.0 (as of 7/21/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "crossplane",
			Name:        "core",
			Description: "Crossplane Core Components",
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "crossplane",
			ReleaseName:   "crossplane",
			RepositoryURL: "https://charts.crossplane.io/stable/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,

		PostInstallResources: []*resource.Resource{
			waitForProviderCRD(),
			providerFamilyAWS(options),
			waitForProviderConfigCRD(),
			defaultProviderConfig(),
		},
	}
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
`
