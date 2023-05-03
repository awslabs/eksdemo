package crossplane

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://crossplane.io/docs/
// GitHub:  https://github.com/crossplane/crossplane
// Helm:    https://github.com/crossplane/crossplane/tree/master/cluster/charts/crossplane
// Repo:    crossplane/crossplane
// Version: Latest is Chart/App v1.9.0, Provider AWS v0.29.0 (as of 08/07/22)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Name:        "crossplane",
			Description: "Cloud Native Control Planes",
		},

		Dependencies: []*resource.Resource{
			crossplaneIrsa(options),
		},

		Options: &application.ApplicationOptions{},

		Installer: &installer.HelmInstaller{
			ChartName:     "crossplane",
			ReleaseName:   "crossplane",
			RepositoryURL: "https://charts.crossplane.io/stable/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		PostInstallResources: []*resource.Resource{
			waitForControllerConfigCRD(),
			awsProvider(options),
			waitForProviderConfigCRD(),
			defaultProviderConfig(),
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
`
