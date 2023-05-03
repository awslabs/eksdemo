package adot_collector

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-EKS-otel.html
// GitHub:  https://github.com/aws-observability/aws-otel-collector
// Helm:    https://github.com/aws-observability/aws-otel-helm-charts/tree/main/charts/adot-exporter-for-eks-on-ec2
// Repo:    amazon/aws-otel-collector
// Version: Latest is Chart 0.6.0, app v0.20.0 (as of 07/30/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "container-insights",
			Name:        "adot-collector",
			Description: "Container Insights ADOT Collector Metrics",
			Aliases:     []string{"adot"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-adot-collector-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "adot-exporter-for-eks-on-ec2",
			ReleaseName:   "container-insights-adot-collector",
			RepositoryURL: "https://aws-observability.github.io/aws-otel-helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace:      "amazon-metrics",
			ServiceAccount: "adot-collector-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.6.0",
				Latest:        "v0.20.0",
				PreviousChart: "0.6.0",
				Previous:      "v0.20.0",
			},
		},
	}

	return app
}

const valuesTemplate = `---
clusterName: {{ .ClusterName }}
fluentbit:
  enabled: false
adotCollector:
  image:
    tag: {{ .Version }}
  daemonSet:
    enabled: true
    createNamespace: false
    serviceAccount:
      create: true
      name: {{ .ServiceAccount }}
      annotations:
        {{ .IrsaAnnotation }}
    service:
      metrics:
        receivers: ["awscontainerinsightreceiver"]
        processors: ["batch/metrics"]
        exporters: ["awsemf"]
`
