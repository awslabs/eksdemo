package falco

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://falco.org/docs/
// GitHub:  https://github.com/falcosecurity/falco
// Helm:    https://github.com/falcosecurity/charts/tree/master/falco
// Repo:    falcosecurity/falco
// Version: Latest is 0.32.0 (as of 06/23/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "falco",
			Description: "Cloud Native Runtime Security",
			Aliases:     []string{"falco-security"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "falco",
			ReleaseName:   "falco",
			RepositoryURL: "https://falcosecurity.github.io/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return addOptions(app)
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
serviceAccount:
  name: {{ .ServiceAccount }}
fakeEventGenerator:
  enabled: {{ .EventGenerator }}
  replicas: {{ .Replicas }}
`
