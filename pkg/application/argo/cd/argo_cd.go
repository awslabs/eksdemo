package cd

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://argo-cd.readthedocs.io/
// GitHub:  https://github.com/argoproj/argo-cd
// Helm:    https://github.com/argoproj/argo-helm/tree/main/charts/argo-cd
// Repo:    quay.io/argoproj/argocd
// Version: Latest Chart is 7.5.2, Argo CD v2.12.3 (as of 9/8/24)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "argo",
			Name:        "cd",
			Description: "Declarative continuous deployment for Kubernetes",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "argo-cd",
			ReleaseName:   "argo-cd",
			RepositoryURL: "https://argoproj.github.io/argo-helm",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
fullnameOverride: argocd
global:
  image:
    tag: {{ .Version }}
configs:
{{- if .IngressHost }}
  params:
    server.insecure: true
{{- end }}
  secret:
    # -- Bcrypt hashed admin password
    argocdServerAdminPassword: "{{ .AdminPassword | bcrypt }}"
server:
  service:
    type: {{ .ServiceType }}
{{- if .IngressHost }}
  ingress:
    enabled: true
  {{- if eq .IngressClass "alb" }}
    controller: aws
    aws:
      serviceType: {{ .ServiceType }}
  {{- end }}
    annotations:
      {{- .IngressAnnotations | nindent 6 }}
    ingressClassName: {{ .IngressClass }}
    hostname: {{ .IngressHost }}
    extraTls:
    - hosts:
      - {{ .IngressHost }}
    {{- if ne .IngressClass "alb" }}
      secretName: argocd-server-tls
    {{- end}}
  {{- if eq .IngressClass "alb" }}
  #ingressGrpc:
  #  enabled: true
  {{- end }}
{{- end }}
`
