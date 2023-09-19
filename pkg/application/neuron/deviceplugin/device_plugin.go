package deviceplugin

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://awsdocs-neuron.readthedocs-hosted.com/en/latest/containers/tutorials/k8s-setup.html
// Repo:    gallery.ecr.aws/neuron/neuron-device-plugin
// Version: Latest is Neuron SDK 2.14.0, Plugin version 2.16.18.0 (as of 9/18/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "neuron-device-plugin",
			Description: "Neuron SDK Device Plugin",
			Aliases:     []string{"ndp"},
		},

		Options: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			Namespace:                 "kube-system",
			ServiceAccount:            "neuron-device-plugin",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "2.16.18.0",
				Previous: "2.16.18.0",
			},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "ai-neuron-device-plugin",
			ResourceTemplate: &template.TextTemplate{
				Template: daemonsetTemplate + rbacTemplate,
			},
		},
	}
	return app
}
