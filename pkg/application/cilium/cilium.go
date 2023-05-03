package cilium

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.cilium.io/
// GitHub:  https://github.com/cilium/cilium
// Helm:    https://github.com/cilium/cilium/tree/master/install/kubernetes/cilium
// Repo:    https://quay.io/repository/cilium/cilium
// Version: Latest is v1.12.6 (as of 01/30/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cilium",
			Description: "eBPF-based Networking, Observability, Security",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cilium",
			ReleaseName:   "cilium",
			RepositoryURL: "https://helm.cilium.io/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
{{- if not .Overlay }}
cni:
  chainingMode: aws-cni
{{- end }}
{{- if .Wireguard }}
encryption:
  # -- Enable transparent network encryption.
  enabled: true
  # -- Encryption method. Can be either ipsec or wireguard.
  type: wireguard
{{- end }}
{{- if not .Overlay }}
enableIPv4Masquerade: false
tunnel: disabled
{{- end }}
{{- if .Wireguard }}
l7Proxy: false
{{- end }}
`
