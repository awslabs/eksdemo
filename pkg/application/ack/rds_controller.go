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
// GitHub:  https://github.com/aws-controllers-k8s/rds-controller
// Helm:    https://github.com/aws-controllers-k8s/rds-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/rds-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/rds-controller
// Version: Latest is v1.4.4 (as of 9/6/24)

func NewRDSController() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "rds-controller",
			Description: "ACK RDS Controller",
			Aliases:     []string{"rds"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-rds-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/rds-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/AmazonRDSFullAccess"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-rds-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.4",
				Latest:        "1.4.4",
				PreviousChart: "1.4.4",
				Previous:      "1.4.4",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-rds-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/rds-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: rdsValuesTemplate,
			},
		},
	}
}

// https://github.com/aws-controllers-k8s/rds-controller/blob/main/helm/values.yaml
const rdsValuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-rds-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
