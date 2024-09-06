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
// GitHub:  https://github.com/aws-controllers-k8s/efs-controller
// Helm:    https://github.com/aws-controllers-k8s/efs-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/efs-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/efs-controller
// Version: Latest is v1.0.0 (as of 9/6/24)

func NewEFSController() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "efs-controller",
			Description: "ACK EFS Controller",
			Aliases:     []string{"efs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-efs-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/efs-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonElasticFileSystemFullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-efs-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.0",
				Latest:        "1.0.0",
				PreviousChart: "1.0.0",
				Previous:      "1.0.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-efs-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/efs-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: efsValuesTemplate,
			},
		},
	}
}

// https://github.com/aws-controllers-k8s/efs-controller/blob/main/helm/values.yaml
const efsValuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-efs-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
