package linkerd_base

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "linkerd",
			Name:        "linkerd-crds",
			Description: "Linkerd Service Mesh Custom Resource Definitions",
		},

                Options: &application.ApplicationOptions{
                        Namespace: "linkerd",
                },

		Installer: &installer.HelmInstaller{
			ChartName:     "linkerd-crds",
			ReleaseName:   "linkerd-crds",
			RepositoryURL: "https://helm.linkerd.io/edge",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()
	return app
}

// https://github.com/linkerd/linkerd2/blob/main/charts/linkerd-control-plane/values.yaml
const valuesTemplate = ``
