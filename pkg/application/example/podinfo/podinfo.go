package podinfo

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/stefanprodan/podinfo/blob/master/README.md
// GitHub:  https://github.com/stefanprodan/podinfo
// Helm:    https://github.com/stefanprodan/podinfo/tree/master/charts/podinfo
// Repo:    ghcr.io/stefanprodan/podinfo
// Version: Latest is Chart/App 6.2.0 (as of 08/25/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "podinfo",
			Description: "Go app w/microservices best practices",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "podinfo",
			ReleaseName:   "example-podinfo",
			RepositoryURL: "https://stefanprodan.github.io/podinfo",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "6.2.0",
				Latest:        "6.2.0",
				PreviousChart: "6.2.0",
				Previous:      "6.2.0",
			},
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "podinfo",
		},
	}

	return app
}

const valuesTemplate = `---
fullnameOverride: podinfo
image:
  tag: {{ .Version }}
service:
  type: {{ .ServiceType }}
  annotations:
    {{- .ServiceAnnotations | nindent 4 }}
{{- if .IngressHost }}
ingress:
  enabled: true
  className: {{ .IngressClass }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  hosts:
  - host: {{ .IngressHost }}
    paths:
    - path: /
      pathType: Prefix
  tls:
  - hosts:
    - {{ .IngressHost }}
  {{- if ne .IngressClass "alb" }}
    secretName: podinfo-tls
  {{- end }}
{{- end }}
`
