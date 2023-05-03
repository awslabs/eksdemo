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
// Repo:    public.ecr.aws/aws-controllers-k8s/prometheusservice-controller
// Version: Latest is v0.1.1 (as of 10/24/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "prometheusservice-controller",
			Description: "ACK Amazon Managed Prometheus Controller",
			Aliases:     []string{"prometheus", "prom", "amp"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-s3-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/prometheusservice-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonPrometheusConsoleFullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-prometheusservice-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.1.1",
				Latest:        "v0.1.1",
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
