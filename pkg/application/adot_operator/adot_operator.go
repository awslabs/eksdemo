package adot_operator

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://opentelemetry.io/docs/
// Docs:    https://github.com/open-telemetry/opentelemetry-operator/blob/main/docs/api.md
// GitHub:  https://github.com/open-telemetry/opentelemetry-operator
// GitHub:  https://github.com/aws-observability/aws-otel-collector/
// Helm:    https://github.com/open-telemetry/opentelemetry-helm-charts/tree/main/charts/opentelemetry-operator
// Repo:    https://gallery.ecr.aws/aws-observability/adot-operator
// Version: Latest is ADOT Collector v0.30.0, Operator v0.78.0, OTEL Chart 0.31.0, (as of 6/18/23)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "adot-operator",
			Description: "AWS Distro for OpenTelemetry (ADOT) Operator",
			Aliases:     []string{"adot"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name:           "adot-collector-irsa",
					ServiceAccount: options.CollectorServiceAccount,
				},
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess",
					"arn:aws:iam::aws:policy/AWSXrayWriteOnlyAccess",
					"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy",
				},
			}),
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "opentelemetry-operator",
			ReleaseName:   "adot-operator",
			RepositoryURL: "https://open-telemetry.github.io/opentelemetry-helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,

		PostInstallResources: []*resource.Resource{
			collectorServiceAccount(options),
		},
	}
}

// This values file is configured to match the ADOT configuration of the EKS Managed Addon
// https://github.com/open-telemetry/opentelemetry-helm-charts/blob/main/charts/opentelemetry-operator/values.yaml
const valuesTemplate = `---
replicaCount: 1
nameOverride: adot-operator
manager:
  image:
    # Images:   gallery.ecr.aws/aws-observability/adot-operator
    repository: public.ecr.aws/aws-observability/adot-operator
    tag: {{ .Version }}
  collectorImage:
    # Images:   gallery.ecr.aws/aws-observability/aws-otel-collector
    repository: public.ecr.aws/aws-observability/aws-otel-collector
    tag: v0.30.0
  targetAllocatorImage:
    # Docs:     https://github.com/open-telemetry/opentelemetry-operator#target-allocator
    # Docs:     https://github.com/open-telemetry/opentelemetry-operator/blob/main/cmd/otel-allocator/README.md
    # Images:   gallery.ecr.aws/aws-observability/adot-operator-targetallocator
    repository: public.ecr.aws/aws-observability/adot-operator-targetallocator
    tag: 0.78.0
  autoInstrumentationImage:
    java:
      # Images:   gallery.ecr.aws/aws-observability/adot-autoinstrumentation-java
      repository: public.ecr.aws/aws-observability/adot-autoinstrumentation-java
      tag: v1.26.0
    nodejs:
      # Images:   gallery.ecr.aws/aws-observability/adot-operator-autoinstrumentation-nodejs
      repository: public.ecr.aws/aws-observability/adot-operator-autoinstrumentation-nodejs
      tag: 0.39.1
    python:
      # Images:   gallery.ecr.aws/aws-observability/adot-operator-autoinstrumentation-python
      repository: public.ecr.aws/aws-observability/adot-operator-autoinstrumentation-python
      tag: 0.39b0
    dotnet:
      # Images:   gallery.ecr.aws/aws-observability/adot-operator-autoinstrumentation-dotnet
      repository: public.ecr.aws/aws-observability/adot-operator-autoinstrumentation-dotnet
      tag: 0.7.0
  serviceAccount:
    name: {{ .ServiceAccount }}
  extraArgs:
    # Images: gallery.ecr.aws/aws-observability/adot-operator-opamp-bridge
    - --operator-opamp-bridge-image=public.ecr.aws/aws-observability/adot-operator-opamp-bridge:0.78.0
kubeRBACProxy:
  image:
    # Images:   gallery.ecr.aws/aws-observability/mirror-kube-rbac-proxy
    repository: public.ecr.aws/aws-observability/mirror-kube-rbac-proxy
    tag: v0.14.1
`
