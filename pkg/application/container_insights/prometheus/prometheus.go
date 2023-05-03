package prometheus

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/ContainerInsights-Prometheus-Setup.html
// GitHub:   https://github.com/aws-samples/amazon-cloudwatch-container-insights, see Releases for versions
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/service/cwagent-prometheus/prometheus-eks.yaml
// Repo:     https://gallery.ecr.aws/cloudwatch-agent/cloudwatch-agent
// Version:  Latest is 1.247352.0b251908 aka k8s/1.3.10 (as of 07/16/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "container-insights",
			Name:        "prometheus",
			Description: "Container Insights monitoring for Prometheus",
			Aliases:     []string{"prom"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-prometheus-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Options: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			Namespace:                 "amazon-cloudwatch",
			ServiceAccount:            "cwagent-prometheus",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.247352.0b251908",
				Previous: "1.247350.0b251780",
			},
		},

		Installer: &installer.ManifestInstaller{
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},
	}

	return app
}
