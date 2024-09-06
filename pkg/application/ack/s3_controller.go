package ack

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
// GitHub:  https://github.com/aws-controllers-k8s/s3-controller
// Helm:    https://github.com/aws-controllers-k8s/s3-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/s3-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/s3-controller
// Version: Latest is v1.0.16 (as of 9/6/24)

func NewS3Controller() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "s3-controller",
			Description: "ACK S3 Controller",
			Aliases:     []string{"s3"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-s3-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/s3-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonS3FullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-s3-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.16",
				Latest:        "1.0.16",
				PreviousChart: "1.0.4",
				Previous:      "1.0.4",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-s3-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/s3-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: s3ValuesTemplate,
			},
		},
	}
}

// https://github.com/aws-controllers-k8s/s3-controller/blob/main/helm/values.yaml
const s3ValuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-s3-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
