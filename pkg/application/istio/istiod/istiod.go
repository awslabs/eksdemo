package istiod

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://istio.io/latest/docs/
// GitHub:  https://github.com/istio/istio
// Helm:    https://github.com/istio/istio/tree/master/manifests/charts/istio-control/istio-discovery
// Repo:    https://hub.docker.com/r/istio/pilot
// Version: Latest is v1.14.1 (as of 06/25/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "istio",
			Name:        "istiod",
			Description: "Istio Control Plane",
			Aliases:     []string{"control-plane", "control", "cp"},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "istio-system",
			ServiceAccount: "istiod",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.14.1",
				Latest:        "1.14.1",
				PreviousChart: "1.14.0",
				Previous:      "1.14.0",
			},
			// Service Account name is hard coded in the Chart
			// https://github.com/istio/istio/blob/master/manifests/charts/istio-control/istio-discovery/templates/serviceaccount.yaml#L10
			DisableServiceAccountFlag: true,
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "istiod",
			ReleaseName:   "istio-istiod",
			RepositoryURL: "https://istio-release.storage.googleapis.com/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		// TODO: option to choose namespace for monitors

		// PostInstallResources: []*resource.Resource{
		// 	podMonitor(),
		// 	serviceMonitor(),
		// },
	}
	return app
}

const valuesTemplate = `---
pilot:
  tag: {{ .Version }}
global:
  istioNamespace: {{ .Namespace }}
`
