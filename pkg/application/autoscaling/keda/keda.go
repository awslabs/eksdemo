package keda

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://keda.sh/docs/
// GitHub:  https://github.com/kedacore/keda
// Helm:    https://github.com/kedacore/charts/tree/main/keda
// Repo:    ghcr.io/kedacore/keda
// Version: Latest is chart 2.7.2, app v2.7.1 (as of 07/26/22)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
			Name:        "keda",
			Description: "Kubernetes-based Event Driven Autoscaling",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "keda",
			ReleaseName:   "autoscaling-keda",
			RepositoryURL: "https://kedacore.github.io/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

const valuesTemplate = `---
serviceAccount:
  name: {{ .ServiceAccount }}
`
