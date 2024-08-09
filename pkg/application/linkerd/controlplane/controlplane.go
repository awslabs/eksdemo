package controlplane

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
	options := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "linkerd",
			Name:        "control-plane",
			Description: "Linkerd Service Mesh Control Plane",
			Aliases:     []string{"controlplane", "cp"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "linkerd-control-plane",
			ReleaseName:   "linkerd-control-plane",
			RepositoryURL: "https://helm.linkerd.io/edge",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}
}

// https://github.com/linkerd/linkerd2/blob/main/charts/linkerd-control-plane/values.yaml
const valuesTemplate = `---
identityTrustAnchorsPEM: |
  {{- .TrustAnchor | trim | nindent 2 }}
identity:
  scheme: linkerd.io/tls
  issuer:
    tls:
      crtPEM: |
        {{- .IssuerCert | trim | nindent 8 }}
      keyPEM: |
        {{- .IssuerKey | trim | nindent 8 }}
`
