package headlamp

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://headlamp.dev/docs/latest/installation/in-cluster/
// GitHub:  https://github.com/headlamp-k8s/headlamp
// Helm:    https://headlamp-k8s.github.io/headlamp/
// Repo:    ghcr.io/headlamp-k8s/headlamp
// Version: Latest is Chart/App 0.15.0/0.20.0 (as of 09/29/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "headlamp",
			Description: "An easy-to-use and extensible Kubernetes web UI",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "headlamp",
			ReleaseName:   "headlamp",
			RepositoryURL: "https://headlamp-k8s.github.io/headlamp/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.15.0",
				Latest:        "0.15.0",
				PreviousChart: "0.15.0",
				Previous:      "0.15.0",
			},
			Namespace: "headlamp",
		},
	}

	return app
}

const valuesTemplate = `---
replicaCount: 1
image:
  tag: "v0.20.0"

config:
  baseURL: ""
  pluginsDir: "/headlamp/plugins"

serviceAccount:
  create: true

clusterRoleBinding:
  create: true

persistentVolumeClaim:
  enabled: false
`
