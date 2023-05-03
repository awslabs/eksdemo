package vpa

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes/autoscaler/blob/master/vertical-pod-autoscaler/README.md
// GitHub:  https://github.com/kubernetes/autoscaler/tree/master/vertical-pod-autoscaler
// Helm:    https://github.com/FairwindsOps/charts/tree/master/stable/vpa
// Repo:    registry.k8s.io/autoscaling/vpa-recommender, registry.k8s.io/autoscaling/vpa-updater
// Version: Latest is chart 1.7.2, VPA 0.13.0 (as of 4/23/23)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
			Name:        "vpa",
			Description: "Vertical Pod Autoscaler",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "vpa",
			ReleaseName:   "autoscaling-vpa",
			RepositoryURL: "https://charts.fairwinds.com/stable",
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
fullnameOverride: vpa
recommender:
  image:
    tag: {{ .Version }}
updater:
  image:
    tag: {{ .Version }}
admissionController:
  enabled: {{ .AdmissionControllerEnabled }}
  image:
    tag: {{ .Version }}
`
