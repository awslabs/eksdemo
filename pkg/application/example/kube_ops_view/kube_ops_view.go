package kube_ops_view

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Codeberg: https://codeberg.org/hjacobs/kube-ops-view
// Manifest: https://codeberg.org/hjacobs/kube-ops-view/src/branch/main/deploy
// Repo:     hjacobs/kube-ops-view

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "kube-ops-view",
			Description: "Kubernetes Operational View",
			Aliases:     []string{"kubeopsview"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-kube-ops-view",
			ResourceTemplate: &template.TextTemplate{
				Template: deploymentTemplate + rbacTemplate + serviceTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app
}
