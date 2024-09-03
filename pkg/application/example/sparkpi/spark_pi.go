package sparkpi

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// GitHub: https://github.com/kubeflow/spark-operator/blob/master/examples/

func NewApp() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "spark-pi",
			Description: "Apache Spark SparkPi Example",
			Aliases:     []string{"spark"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-spark-pi",
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace:          "default",
			DisableVersionFlag: true,
			ServiceAccount:     "spark-operator-spark",
		},
	}
}
