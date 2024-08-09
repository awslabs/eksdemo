package crds

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "linkerd",
			Name:        "crds",
			Description: "Linkerd Service Mesh Custom Resource Definitions",
			Aliases:     []string{"crd"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "linkerd-crds",
			ReleaseName:   "linkerd-crds",
			RepositoryURL: "https://helm.linkerd.io/edge",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2024.7.3",
				PreviousChart: "2024.7.3",
			},
			Namespace: "linkerd",
		},
	}
}

// https://github.com/linkerd/linkerd2/blob/main/charts/linkerd-crds/values.yaml
const valuesTemplate = ``
