package metrics_server

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/metrics-server/blob/master/README.md
// GitHub:  https://github.com/kubernetes-sigs/metrics-server
// Helm:    https://github.com/kubernetes-sigs/metrics-server/tree/master/charts/metrics-server
// Repo:    registry.k8s.io/metrics-server/metrics-server
// Version: Latest is Chart 3.10.0, App v0.6.3 (as of 4/23/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "metrics-server",
			Description: "Kubernetes Metric Server",
			Aliases:     []string{"metrics"},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "metrics-server",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "3.10.0",
				Latest:        "v0.6.3",
				PreviousChart: "3.8.3",
				Previous:      "v0.6.2",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "metrics-server",
			ReleaseName:   "metrics-server",
			RepositoryURL: "https://kubernetes-sigs.github.io/metrics-server/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
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
