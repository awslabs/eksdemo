package prometheusservice_controller

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://aws-controllers-k8s.github.io/community/docs/community/overview/
// Docs:    https://aws-controllers-k8s.github.io/community/reference/
// GitHub:  https://github.com/aws-controllers-k8s/prometheusservice-controller
// Helm:    https://github.com/aws-controllers-k8s/prometheusservice-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/prometheusservice-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/prometheusservice-controller
// Version: Latest is v1.2.3 (as of 6/11/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "prometheusservice-controller",
			Description: "ACK Prometheus Service Controller",
			Aliases:     []string{"prometheus", "prom", "amp"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-prometheusservice-controller-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocTemplate,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-prometheusservice-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.2.3",
				Latest:        "1.2.3",
				PreviousChart: "v0.1.1",
				Previous:      "v0.1.1",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-prometheusservice-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/prometheusservice-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

// https://github.com/aws-controllers-k8s/prometheusservice-controller/blob/main/config/iam/recommended-inline-policy
const policyDocTemplate = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - aps:*
  - logs:CreateLogDelivery
  - logs:DescribeLogGroups
  - logs:DescribeResourcePolicies
  - logs:PutResourcePolicy
  Resource: "*"
`

const valuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-prometheusservice-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
