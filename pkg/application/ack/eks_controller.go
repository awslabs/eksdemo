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
// GitHub:  https://github.com/aws-controllers-k8s/eks-controller
// Helm:    https://github.com/aws-controllers-k8s/eks-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/eks-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/eks-controller
// Version: Latest is v1.4.5 (as of 9/5/24)

func NewEKSController() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "eks-controller",
			Description: "ACK EKS Controller",
			Aliases:     []string{"eks"},
		},

		Dependencies: []*resource.Resource{
			fargatePodExecutionRole(),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-eks-controller-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: eksPolicyDocTemplate,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-eks-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.4.5",
				Latest:        "1.4.5",
				PreviousChart: "1.0.2",
				Previous:      "1.0.2",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-eks-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/eks-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: eksValuesTemplate,
			},
		},
	}
}

// https://github.com/aws-controllers-k8s/eks-controller/blob/main/config/iam/recommended-inline-policy
const eksPolicyDocTemplate = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - eks:*
  Resource: "*"
- Effect: Allow
  Action:
  - iam:GetRole
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:role/aws-service-role/eks-fargate.amazonaws.com/AWSServiceRoleForAmazonEKSForFargate
- Effect: Allow
  Action:
  - iam:PassRole
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:role/eksdemo.{{ .ClusterName }}.fargate-pod-execution-role
`

const eksValuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-eks-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
