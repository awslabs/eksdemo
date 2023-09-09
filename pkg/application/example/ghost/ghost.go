package ghost

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://ghost.org/docs/
// GitHub:  https://github.com/TryGhost/Ghost
// Helm:    https://github.com/bitnami/charts/tree/main/bitnami/ghost
// Repo:    https://hub.docker.com/r/bitnami/ghost
// Version: Latest is Chart 19.5.5, App 5.62.0 (as of 9/12/23)

func New() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "ghost",
			Description: "Turn your audience into a business",
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "ghost",
			ReleaseName:   ghostReleaseName,
			RepositoryURL: "oci://registry-1.docker.io/bitnamicharts/ghost",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PVCLabels: map[string]string{
				"app.kubernetes.io/instance": ghostReleaseName,
			},
		},

		Options: options,
	}
}

const ghostReleaseName = `example-ghost`

const valuesTemplate = `---
{{- if .StorageClass }}
global:
  storageClass: {{ .StorageClass }}
{{- end }}
fullnameOverride: ghost
image:
  tag: {{ .Version }}
ghostPassword: {{ .GhostPassword }}
ghostHost: {{ .IngressHost }}
ghostEnableHttps: true
service:
  type: {{ .ServiceType }}
  annotations:
    {{- .ServiceAnnotations | nindent 4 }}
ingress:
  enabled: true
  pathType: Prefix
  hostname: {{ .IngressHost }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  tls: true
  ingressClassName: {{ .IngressClass }}
mysql:
  fullnameOverride: ghost-mysql
`
