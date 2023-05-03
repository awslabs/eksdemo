package flux_controllers

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://fluxcd.io/docs/
// GitHub:  https://github.com/fluxcd/flux2/
// Helm:    https://github.com/fluxcd-community/helm-charts/tree/main/charts/flux2
// Repo:    ghcr.io/fluxcd/<resource>-controller
// Version: Latest Chart is 1.0.0, Flux v0.31.3 (as of 07/09/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "flux",
			Name:        "controllers",
			Description: "Flux Controllers",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "flux2",
			ReleaseName:   "flux2",
			RepositoryURL: "https://fluxcd-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.0",
				Latest:        "v0.31.3",
				PreviousChart: "0.20.0",
				Previous:      "v0.31.1",
			},
			DisableServiceAccountFlag: true,
			LockVersionFlag:           true,
			Namespace:                 "flux-system",
		},
	}

	return app
}

const valuesTemplate = ``
