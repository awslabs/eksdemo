package contour

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://projectcontour.io/docs/
// GitHub:  https://github.com/projectcontour/contour
// Helm:    https://github.com/bitnami/charts/tree/master/bitnami/contour
// Repo:    docker.io/bitnami/contour, docker.io/bitnami/envoy
// Version: Latest is Chart 9.0.3, Contour v1.22.0 (as of 08/20/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ingress",
			Name:        "contour",
			Description: "Ingress Controller using Envoy proxy",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "contour",
			ReleaseName:   "ingress-contour",
			RepositoryURL: "https://charts.bitnami.com/bitnami",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
fullnameOverride: contour
contour:
  image:
    tag: {{ .Version }}
  replicaCount: {{ .Replicas }}
envoy:
  service:
    annotations: {}
`
