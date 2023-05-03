package ec2_controller

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
// GitHub:  https://github.com/aws-controllers-k8s/ec2-controller
// Helm:    https://github.com/aws-controllers-k8s/ec2-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/ec2-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/ec2-controller
// Version: Latest is v0.0.21 (as of 10/24/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "ec2-controller",
			Description: "ACK EC2 Controller",
			Aliases:     []string{"ec2"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-ec2-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/ec2-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonEC2FullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-ec2-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.0.21",
				Latest:        "v0.0.21",
				PreviousChart: "v0.0.15",
				Previous:      "v0.0.15",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-ec2-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/ec2-chart",
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
fullnameOverride: ack-ec2-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
