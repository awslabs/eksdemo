package prometheus_node_exporter

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://prometheus.io/docs/guides/node-exporter/
// GitHub:  https://github.com/prometheus/node_exporter
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/prometheus-node-exporter
// Repo:    quay.io/prometheus/node-exporter
// Version: Latest is v1.5.0 (as of 2/8/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "prometheus-node-exporter",
			Description: "Prometheus Node Exporter",
			Aliases:     []string{"node-exporter"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "prometheus-node-exporter",
			ReleaseName:   "prometheus-node-exporter",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "monitoring",
			ServiceAccount: "prometheus-node-exporter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "4.13.0",
				Latest:        "v1.5.0",
				PreviousChart: "4.13.0",
				Previous:      "v1.5.0",
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
rbac:
  pspEnabled: false
`
