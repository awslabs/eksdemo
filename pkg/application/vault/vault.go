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
  # For HA configuration and because we need to manually init the vault,
  # we need to define custom readiness/liveness Probe settings
    raft:
       enabled: true
       setNodeId: true
       config: |
          cluster_name = "vault-integrated-storage"
          storage "raft" {
             path    = "/vault/data/"
          }
          listener "tcp" {
             address = "[::]:8200"
             cluster_address = "[::]:8201"
             tls_disable = "true"
          }
          service_registration "kubernetes" {}
{{ else }}
  ha:
    enabled: false
    replicas: 1
{{ end }}`
