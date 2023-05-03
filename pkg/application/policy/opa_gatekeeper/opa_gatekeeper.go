package opa_gatekeeper

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://open-policy-agent.github.io/gatekeeper/website/docs/
// GitHub:  https://github.com/open-policy-agent/gatekeeper
// Helm:    https://github.com/open-policy-agent/gatekeeper/tree/master/charts/gatekeeper
// Repo:    openpolicyagent/gatekeeper
// Version: Latest is v3.9.0 (as of 07/22/22)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "policy",
			Name:        "opa-gatekeeper",
			Description: "Policy Controller for Kubernetes",
			Aliases:     []string{"opa", "gatekeeper"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "gatekeeper",
			ReleaseName:   "policy-opa-gatekeeper",
			RepositoryURL: "https://open-policy-agent.github.io/gatekeeper/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

const valuesTemplate = `---
replicas: {{ .Replicas }}
image:
  release: {{ .Version }}
psp:
  enabled: false
`
