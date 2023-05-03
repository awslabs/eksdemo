package kube_state_metrics

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes/kube-state-metrics/tree/main/docs
// GitHub:  https://github.com/kubernetes/kube-state-metrics
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-state-metrics
// Repo:    registry.k8s.io/kube-state-metrics/kube-state-metrics
// Version: Latest is Chart 4.29.0, App v2.7.0 (as of 2/13/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "kube-state-metrics",
			Description: "Kube State Metrics",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kube-state-metrics",
			ReleaseName:   "kube-state-metrics",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "kube-state-metrics",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.29.0",
				Latest:        "v2.7.0",
				PreviousChart: "4.29.0",
				Previous:      "v2.7.0",
			},
		},
	}

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
serviceAccount:
  name: {{ .ServiceAccount }}
`
