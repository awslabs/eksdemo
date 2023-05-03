package adot_operator

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://opentelemetry.io/docs/
// Docs:    https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md
// GitHub:  https://github.com/open-telemetry/opentelemetry-operator
// Helm:    https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator
// Repo:    https://gallery.ecr.aws/aws-observability/adot-operator
// Repo:    https://gallery.ecr.aws/aws-observability/mirror-kube-rbac-proxy
// Version: Latest is v0.66.0, Chart 0.21.4 (as of 01/24/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "adot-operator",
			Description: "AWS Distro for OpenTelemetry (ADOT) Operator",
			Aliases:     []string{"adot"},
		},

		Options: &application.ApplicationOptions{
			Namespace: "adot-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.21.4",
				Latest:        "v0.66.0",
				PreviousChart: "0.17.2",
				Previous:      "v0.63.1",
			},
			// Service Account name isn't flexible
			// https://github.com/open-telemetry/opentelemetry-helm-charts/blob/main/charts/opentelemetry-operator/templates/serviceaccount.yaml
			DisableServiceAccountFlag: true,
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "opentelemetry-operator",
			ReleaseName:   "adot-operator",
			RepositoryURL: "https://open-telemetry.github.io/opentelemetry-helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
replicaCount: 1
nameOverride: adot-operator
manager:
  image:
    repository: public.ecr.aws/aws-observability/adot-operator
    tag: {{ .Version }}
  collectorImage:
    repository: public.ecr.aws/aws-observability/aws-otel-collector
    tag: v0.24.1
  targetAllocatorImage:
    repository: public.ecr.aws/aws-observability/mirror-target-allocator
    tag: 0.66.0
  autoInstrumentationImage:
    java:
      repository: public.ecr.aws/aws-observability/mirror-autoinstrumentation-java
      tag: 1.20.2
    nodejs:
      repository: public.ecr.aws/aws-observability/mirror-autoinstrumentation-nodejs
      tag: 0.31.0
    python:
      repository: public.ecr.aws/aws-observability/mirror-autoinstrumentation-python
      tag: 0.35b0
    dotnet:
      repository: public.ecr.aws/aws-observability/mirror-autoinstrumentation-dotnet
      tag: 0.5.0
kubeRBACProxy:
  image:
    repository: public.ecr.aws/aws-observability/mirror-kube-rbac-proxy
    tag: v0.11.0
`
