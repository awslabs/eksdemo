package kyverno

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://kyverno.io/docs/
// GitHub:  https://github.com/kyverno/kyverno/
// Helm:    https://github.com/kyverno/kyverno/tree/main/charts/kyverno
// Repo:    ghcr.io/kyverno/kyverno
// Version: Latest is Chart 2.5.2, App v1.7.2 (as of 08/08/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "policy",
			Name:        "kyverno",
			Description: "Kubernetes Native Policy Management",
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kyverno",
			ServiceAccount: "kyverno",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.5.2",
				Latest:        "v1.7.2",
				PreviousChart: "2.5.2",
				Previous:      "v1.7.2",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kyverno",
			ReleaseName:   "policy-kyverno",
			RepositoryURL: "https://kyverno.github.io/kyverno/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	return app
}

const valuesTemplate = `---
fullnameOverride: kyverno
rbac:
  serviceAccount:
    name: {{ .ServiceAccount }}
image:
  tag: {{ .Version }}
`
