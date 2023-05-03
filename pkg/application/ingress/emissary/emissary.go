package emissary

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://www.getambassador.io/docs/emissary/
// GitHub:  https://github.com/emissary-ingress/emissary
// Helm:    https://github.com/emissary-ingress/emissary/tree/master/charts/emissary-ingress
// Repo:    docker.io/emissaryingress/emissary
// Version: Latest is Chart 8.0.0 and App 3.0.0 (as of 07/10/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ingress",
			Name:        "emissary",
			Description: "Open Source API Gateway from Ambassador",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "emissary-ingress",
			ReleaseName:   "ingress-emissary",
			RepositoryURL: "https://app.getambassador.io",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

// TODO: Consider proxy protocol if in the future the Helm chart can configure it
//       https://github.com/emissary-ingress/emissary/issues/3300

const valuesTemplate = `---
fullnameOverride: emissary-ingress
replicaCount: {{ .Replicas }}
image:
  tag: {{ .Version }}
service:
  externalTrafficPolicy: Local
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
serviceAccount:
  name: {{ .ServiceAccount }}
createDefaultListeners: true
`
