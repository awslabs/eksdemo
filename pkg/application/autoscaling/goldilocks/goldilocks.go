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
// Version: Latest is chart 6.1.4, app v4.3.3 (as of 07/26/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
			Name:        "goldilocks",
			Description: "Get your resource requests \"Just Right\"",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "goldilocks",
			ReleaseName:   "autoscaling-goldilocks",
			RepositoryURL: "https://charts.fairwinds.com/stable",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "6.1.4",
				Latest:        "v4.3.3",
				PreviousChart: "6.1.4",
				Previous:      "v4.3.3",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "goldilocks",
		},
	}

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: goldilocks
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
