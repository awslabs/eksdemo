package flux_sync

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://fluxcd.io/docs/
// GitHub:  https://github.com/fluxcd/flux2/
// Helm:    https://github.com/fluxcd-community/helm-charts/tree/main/charts/flux2-sync
// Version: Latest Chart is 1.0.0 (as of 07/09/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "flux",
			Name:        "sync",
			Description: "Flux GitRepository to sync with",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "flux2-sync",
			ReleaseName:   "flux-sync",
			RepositoryURL: "https://fluxcd-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = NewOptions()

	return app
}

const valuesTemplate = `---
gitRepository:
  spec:
    # -- The repository URL, can be an HTTP/S or SSH address.
    url: {{ .GitUrl }}
kustomization:
  spec:
    # -- _Optional_ Path to the directory containing the kustomization.yaml
    # file, or the set of plain YAMLs a kustomization.yaml should
    # be generated for. Defaults to ‘None’, which translates to
    # the root path of the SourceRef.
    path: {{ .KustomizationPath }}
	# -- _Optional_ TargetNamespace sets or overrides the namespace in the kustomization.yaml file.
    targetNamespace: {{ .TargetNamespace }}
`
