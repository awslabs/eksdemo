package controlPlane

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
/*
Begin Helm App
*/

	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "linkerd",
			Name:        "linkerd-control-plane",
			Description: "Linkerd Service Mesh Custom Resource Definitions",
		},
                Options: options,
		Installer: &installer.HelmInstaller{
			ChartName:     "linkerd-control-plane",
			ReleaseName:   "linkerd-control-plane",
			RepositoryURL: "https://helm.linkerd.io/edge",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

        app.Flags = flags

	return app
}

// https://github.com/linkerd/linkerd2/blob/main/charts/linkerd-control-plane/values.yaml
const valuesTemplate = `---
identityTrustAnchorsPEM: |
  {{ .CAPEM }}
identity:
    scheme: linkerd.io/tls
    issuer:
        tls:
            crtPEM: |
              {{ .CRTPEM }}
            keyPEM: |
              {{ .KEYPEM }}
`
