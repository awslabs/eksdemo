package harbor

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://goharbor.io/docs
// GitHub:  https://github.com/goharbor/harbor
// Helm:    https://github.com/goharbor/harbor-helm
// Repo:    goharbor/harbor-core, goharbor/harbor-portal, goharbor/registry-photon
// Version: Latest is Chart 1.9.3, App v2.5.3 (as of 08/14/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "harbor",
			Description: "Cloud Native Registry",
		},

		Options: &application.ApplicationOptions{},

		Installer: &installer.HelmInstaller{
			ChartName:     "harbor",
			ReleaseName:   "harbor",
			RepositoryURL: "https://helm.goharbor.io",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
expose:
  type: ingress
  ingress:
    hosts:
      core: {{ .IngressHost }}
    {{- if .NotaryEnabled }}
      notary: {{ .NotaryHost }}
    {{- end }}
    className: {{ .IngressClass }}
    annotations:
      {{- .IngressAnnotations | nindent 6 }}
externalURL: https://{{ .IngressHost }}
harborAdminPassword: {{ .AdminPassword }}
portal:
  image:
    tag: {{ .Version }}
core:
  image:
    tag: {{ .Version }}
jobservice:
  image:
    tag: {{ .Version }}
registry:
  registry:
    image:
      tag: {{ .Version }}
  controller:
    image:
      tag: {{ .Version }}
chartmuseum:
  image:
    tag: {{ .Version }}
trivy:
  image:
    tag: {{ .Version }}
notary:
  enabled: {{ .NotaryEnabled }}
{{- if .NotaryEnabled }}
  server:
    image:
      tag: {{ .Version }}
  signer:
    image:
      tag: {{ .Version }}
{{- end }}
database:
  internal:
    image:
      tag: {{ .Version }}
redis:
  internal:
    image:
      tag: {{ .Version }}
exporter:
  image:
    tag: {{ .Version }}
`
