package vault

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide
// https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-amazon-eks

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
  # Enables TLS across the cluster to verify authenticity of the Vault servers and clients.
  tls:
    enabled: true

# Configures High Availability Mode for Vault Server
{{ if gt .Replicas 1 }}
server:
  ha:
    enabled: true
    # The number of server agents to run. This determines the fault tolerance of the cluster.
    replicas: {{ .Replicas }}
{{ end }}

# Contains values that configure the Vault UI.
{{ if .EnableUI }}
# Vault UI
ui:
  enabled: true
  serviceType: "LoadBalancer"
  serviceNodePort: null
  externalPort: 8200
{{ end }}`
