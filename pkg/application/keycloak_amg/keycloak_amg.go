package keycloak_amg

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/amg_workspace"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://www.keycloak.org/documentation
// GitHub:  https://github.com/keycloak/keycloak
// GitHub:  https://github.com/bitnami/bitnami-docker-keycloak
// Helm:    https://github.com/bitnami/charts/tree/master/bitnami/keycloak
// Repo:    https://hub.docker.com/r/bitnami/keycloak
// Version: Latest is Chart 9.6.8, App 18.0.2 (as of 08/14/22)

func NewApp() *application.Application {
	options, flags := NewOptions()

	options.AmgOptions = &amg_workspace.AmgOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-grafana",
		},
		Auth: []string{"SAML"},
	}

	app := &application.Application{
		Command: cmd.Command{
			Name:        "keycloak-amg",
			Description: "Keycloak SAML iDP for Amazon Managed Grafana",
			Aliases:     []string{"keycloakamg"},
		},

		Dependencies: []*resource.Resource{
			amg_workspace.NewResourceWithOptions(options.AmgOptions),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "keycloak",
			ReleaseName:   keycloakReleaseName,
			RepositoryURL: "https://charts.bitnami.com/bitnami",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PVCLabels: map[string]string{
				"app.kubernetes.io/instance": keycloakReleaseName,
			},
		},
	}

	app.Options = options
	app.Flags = flags

	return app
}

const keycloakReleaseName = `keycloak-amg`

const valuesTemplate = `---
fullnameOverride: keycloak
image:
  tag: {{ .Version }}
auth:
  adminUser: admin
  adminPassword: {{ .AdminPassword }}
service:
  type: {{ .ServiceType }}
ingress:
  enabled: true
  ingressClassName: {{ .IngressClass }}
  pathType: Prefix
  hostname: {{ .IngressHost }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  tls: true
keycloakConfigCli:
  enabled: true
  command:
  - java
  - -jar
  - /opt/bitnami/keycloak-config-cli/keycloak-config-cli-{{ .Version }}.jar
  configuration:
    eksdemo.json: |
      {
        "realm": "eksdemo",
        "enabled": true,
        "roles": {
          "realm": [
            {
              "name": "admin"
            }
          ]
        },
        "users": [
          {
            "username": "admin",
            "email": "admin@eksdemo",
            "enabled": true,
            "firstName": "Admin",
            "realmRoles": [
              "admin"
            ],
            "credentials": [
              {
                "type": "password",
                "value": "{{ .AdminPassword }}"
              }
            ]
          }
        ],
        "clients": [
          {
            "clientId": "https://{{ .AmgWorkspaceUrl }}/saml/metadata",
            "name": "amazon-managed-grafana",
            "enabled": true,
            "protocol": "saml",
            "adminUrl": "https://{{ .AmgWorkspaceUrl }}/login/saml",
            "redirectUris": [
              "https://{{ .AmgWorkspaceUrl }}/saml/acs"
            ],
            "attributes": {
              "saml.authnstatement": "true",
              "saml.server.signature": "true",
              "saml_name_id_format": "email",
              "saml_force_name_id_format": "true",
              "saml.assertion.signature": "true",
              "saml.client.signature": "false"
            },
            "defaultClientScopes": [],
            "protocolMappers": [
              {
                "name": "name",
                "protocol": "saml",
                "protocolMapper": "saml-user-property-mapper",
                "consentRequired": false,
                "config": {
                  "attribute.nameformat": "Unspecified",
                  "user.attribute": "firstName",
                  "attribute.name": "displayName"
                }
              },
              {
                "name": "email",
                "protocol": "saml",
                "protocolMapper": "saml-user-property-mapper",
                "consentRequired": false,
                "config": {
                  "attribute.nameformat": "Unspecified",
                  "user.attribute": "email",
                  "attribute.name": "mail"
                }
              },
              {
                "name": "role list",
                "protocol": "saml",
                "protocolMapper": "saml-role-list-mapper",
                "config": {
                  "single": "true",
                  "attribute.nameformat": "Unspecified",
                  "attribute.name": "role"
                }
              }
            ]
          }
        ]
      }
`
