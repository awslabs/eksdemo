package spark

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://www.kubeflow.org/docs/components/spark-operator/
// GitHub:  https://github.com/kubeflow/spark-operator
// Helm:    https://github.com/kubeflow/spark-operator/tree/master/charts/spark-operator-chart
// Repo:    https://hub.docker.com/r/kubeflow/spark-operator/
// Version: Latest is v2.0.0-rc.0 (as of 9/2/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "spark-operator",
			Description: "Kubeflow Spark Operator",
			Aliases:     []string{"spark"},
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "spark-operator",
			ReleaseName:   "spark-operator",
			RepositoryURL: "https://kubeflow.github.io/spark-operator",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}
}

// https://github.com/kubeflow/spark-operator/blob/master/charts/spark-operator-chart/values.yaml
const valuesTemplate = `---
image:
  tag: {{ .Version }}
controller:
  serviceAccount:
    name: {{ .ServiceAccount }}
`
