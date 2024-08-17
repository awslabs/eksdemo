package goldilocks

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://goldilocks.docs.fairwinds.com/
// GitHub:  https://github.com/FairwindsOps/goldilocks
// Helm:    https://github.com/FairwindsOps/charts/tree/master/stable/goldilocks
// Repo:    us-docker.pkg.dev/fairwinds-ops/oss/goldilocks
// Version: Latest is chart 8.0.2, app v4.13.0 (as of 8/17/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "goldilocks",
			Description: "Get your resource requests \"Just Right\"",
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "goldilocks",
			ReleaseName:   "goldilocks",
			RepositoryURL: "https://charts.fairwinds.com/stable",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}
}

// https://github.com/FairwindsOps/charts/blob/master/stable/goldilocks/values.yaml
const valuesTemplate = `---
vpa:
  enabled: {{ not .NoVPA }}
image:
  tag: {{ .Version }}
dashboard:
  replicaCount: 1
  service:
    type: {{ .ServiceType }}
    annotations:
      {{- .ServiceAnnotations | nindent 6 }}
{{- if .IngressHost }}
  ingress:
    enabled: true
    ingressClassName: {{ .IngressClass }}
    annotations:
      {{- .IngressAnnotations | nindent 6 }}
    hosts:
    - host: {{ .IngressHost }}
      paths:
      - path: /
        type: Prefix
    tls:
    - hosts:
      - {{ .IngressHost }}
    {{- if ne .IngressClass "alb" }}
      secretName: goldilocks-dashboard-tls
    {{- end }}
{{- end }}
`
