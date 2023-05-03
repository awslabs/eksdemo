package openebs

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://openebs.io/docs/
// GitHub:  https://github.com/openebs/openebs
// Helm:    https://github.com/openebs/charts/tree/main/charts/openebs
// Helm:    https://github.com/openebs/dynamic-localpv-provisioner/tree/develop/deploy/helm/charts
// Helm:    https://github.com/openebs/node-disk-manager/tree/develop/deploy/helm/charts
// Repo:    openebs/provisioner-localpv, openebs/node-disk-manager, openebs/node-disk-operator
// Version: Latest is Chart 3.3.0 (as of 07/30/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "storage",
			Name:        "openebs",
			Description: "Kubernetes storage simplified",
		},

		Options: &application.ApplicationOptions{
			Namespace: "openebs",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "3.3.0",
				PreviousChart: "3.3.0",
			},
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "openebs",
			ReleaseName:   "storage-openebs",
			RepositoryURL: "https://openebs.github.io/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	return app
}

const valuesTemplate = `---
# Disable the embedded RBAC resources as they aren't needed
rbac:
  create: false
serviceAccount:
  create: false
# Disable the embedded localprovisioner templates
localprovisioner:
  enabled: false
# Disable the embedded ndm templates
ndm:
  enabled: false
ndmOperator:
  enabled: false
# Enable the official openebs-ndm chart as a dependency
openebs-ndm:
  enabled: true
  fullnameOverride: openebs-ndm
  ndmOperator:
    fullnameOverride: openebs-ndm-operator
# Enable the official localprovisioner chart as a dependency
localpv-provisioner:
  enabled: true
  fullnameOverride: localpv-provisioner
`
