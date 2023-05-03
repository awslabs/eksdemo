package cloudwatch_agent

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-metrics.html
// GitHub:   https://github.com/aws-samples/amazon-cloudwatch-container-insights, see Releases for versions
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/tree/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cwagent
// Repo:     https://gallery.ecr.aws/cloudwatch-agent/cloudwatch-agent
// Version:  Latest is 1.247352.0b251908 aka k8s/1.3.10 (as of 07/16/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "container-insights",
			Name:        "cloudwatch-agent",
			Description: "Container Insights CloudWatch Agent Metrics",
			Aliases:     []string{"cloudwatch", "cwa"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-cloudwatch-agent-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "container-insights-cloudwatch-agent",
			ResourceTemplate: &template.TextTemplate{
				Template: serviceAccountTemplate + configMapTemplate + daemonSetTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			Namespace:                 "amazon-cloudwatch",
			ServiceAccount:            "cloudwatch-agent",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "1.247352.0b251908",
				Previous: "1.247350.0b251780",
			},
		},
	}

	return app
}
