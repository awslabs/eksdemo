package fluent_bit

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:     https://docs.aws.amazon.com/AmazonCloudWatch/latest/monitoring/Container-Insights-setup-logs-FluentBit.html
// GitHub:   https://github.com/aws/aws-for-fluent-bit
// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml
// Repo:     https://gallery.ecr.aws/aws-observability/aws-for-fluent-bit
// Version:  CloudWatch documentation uses "stable" tag

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "container-insights",
			Name:        "fluent-bit",
			Description: "Container Insights Fluent Bit Logs",
			Aliases:     []string{"fluentbit"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "container-insights-fluent-bit-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"},
			}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "container-insights-fluent-bit",
			ResourceTemplate: &template.TextTemplate{
				Template: configMapTemplate + fluentBitManifestTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}
