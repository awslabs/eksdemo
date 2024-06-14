package vault

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

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
// https://github.com/hashicorp/vault-helm/blob/main/values.yaml
}

const valuesTemplate = `---
server:
{{ if gt .Replicas 1 }}
  ha:
    enabled: true
    replicas: {{ .Replicas }}
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
