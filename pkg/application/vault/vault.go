package vault

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide
// https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-amazon-eks

// fmt.Println(getField(&v, "X"))

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "vault",
			Description: "HashiCorp Vault Secrets and Encryption Management System",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "vault",
			ReleaseName:   "vault",
			RepositoryURL: "https://helm.releases.hashicorp.com",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app

}

const valuesTemplate = `---
global:
  # The main enabled/disabled setting.
  # If true, servers, clients, Vault UI will be enabled.
  enabled: true
  # The prefix used for all resources created in the Helm chart.
  name: null
  # TLS for end-to-end encrypted transport
{{ if .EnableTLS }}
  tlsDisable: false
{{ else }}
  tlsDisable: true
{{ end }}

server:
# Configures High Availability Mode for Vault Server
{{ if gt .Replicas 1 }}
  ha:
    enabled: true
    # The number of server agents to run. This determines the fault tolerance of the cluster.
    replicas: {{ .Replicas }}
{{ end }}
{{ if .Enterprise }}
  # [Enterprise Only] This value refers to a Kubernetes secret that you have
  # created that contains your enterprise license. If you are not using an
  # enterprise image or if you plan to introduce the license key via another
  # route, then leave secretName blank ("") or set it to null.
  # Requires Vault Enterprise 1.8 or later.
  enterpriseLicense:
    # The name of the Kubernetes secret that holds the enterprise license. The
    # secret must be in the same namespace that Vault is installed into.
    secretName: "{{ .Enterprise }}"
    # The key within the Kubernetes secret that holds the enterprise license.
    secretKey: "license"
{{ end }}
# Warning: Vault cannot run in both HA and Development Modes; Development Wins
{{ if .DevelopmentMode }}
  dev:
    enabled: true
{{ end }}
`
