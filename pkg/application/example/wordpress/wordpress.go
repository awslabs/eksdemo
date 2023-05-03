package wordpress

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://codex.wordpress.org/Main_Page
// GitHub:  https://github.com/WordPress/WordPress
// Helm:    https://github.com/bitnami/charts/tree/master/bitnami/wordpress
// Repo:    https://hub.docker.com/r/bitnami/wordpress
// Version: Latest is Chart 15.0.13, App 6.0.1 (as of 08/03/22)

func NewApp() *application.Application {
	options, flags := NewOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "wordpress",
			Description: "WordPress Blog",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "wordpress",
			ReleaseName:   wordpressReleaseName,
			RepositoryURL: "https://charts.bitnami.com/bitnami",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PVCLabels: map[string]string{
				"app.kubernetes.io/instance": wordpressReleaseName,
			},
		},
	}

	app.Options = options
	app.Flags = flags

	return app
}

const wordpressReleaseName = `example-wordpress`

const valuesTemplate = `---
{{- if .StorageClass }}
global:
  storageClass: {{ .StorageClass }}
{{- end }}
fullnameOverride: wordpress
image:
  tag: {{ .Version }}
wordpressPassword: {{ .WordpressPassword }}
service:
  type: {{ .ServiceType }}
  annotations:
    {{- .ServiceAnnotations | nindent 4 }}
{{- if .IngressHost }}
ingress:
  enabled: true
  pathType: Prefix
  ingressClassName: {{ .IngressClass }}
  hostname: {{ .IngressHost }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  tls: true
{{- end }}
mariadb:
  fullnameOverride: wordpress-mariadb
`
