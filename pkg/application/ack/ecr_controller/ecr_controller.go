package ecr_controller

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
// GitHub:  https://github.com/aws-controllers-k8s/ecr-controller
// Helm:    https://github.com/aws-controllers-k8s/ecr-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/ecr-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/ecr-controller
// Version: Latest is v1.0.4 (as of 6/11/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "ecr-controller",
			Description: "ACK ECR Controller",
			Aliases:     []string{"ecr"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-ecr-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/ecr-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryFullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-ecr-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.4",
				Latest:        "1.0.4",
				PreviousChart: "v0.1.7",
				Previous:      "v0.1.7",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-ecr-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/ecr-chart",
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
fullnameOverride: ack-ecr-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
