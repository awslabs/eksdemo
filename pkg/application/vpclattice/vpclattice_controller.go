package vpclattice

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/aws/aws-application-networking-k8s/tree/main/docs
// GitHub:  https://github.com/aws/aws-application-networking-k8s
// Helm:    https://github.com/aws/aws-application-networking-k8s/tree/main/helm
// Repo:    https://gallery.ecr.aws/aws-application-networking-k8s/aws-gateway-controller
// Version: Latest is v1.0.6 (as of 7/25/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "vpc-lattice-controller",
			Description: "Amazon VPC Lattice (Gateway API) Controller",
			Aliases:     []string{"gateway-api-controller", "vpc-lattice", "vpclattice", "lattice"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "vpc-lattice-controller-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
			securityGroupRule(options),
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ReleaseName:   "vpc-lattice-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-application-networking-k8s/aws-gateway-controller-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,

		PostInstallResources: []*resource.Resource{
			gatewayClass(),
		},
	}
}

// https://github.com/aws/aws-application-networking-k8s/blob/main/config/iam/recommended-inline-policy.json
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - vpc-lattice:*
  - ec2:DescribeVpcs
  - ec2:DescribeSubnets
  - ec2:DescribeTags
  - ec2:DescribeSecurityGroups
  - logs:CreateLogDelivery
  - logs:GetLogDelivery
  - logs:DescribeLogGroups
  - logs:PutResourcePolicy
  - logs:DescribeResourcePolicies
  - logs:UpdateLogDelivery
  - logs:DeleteLogDelivery
  - logs:ListLogDeliveries
  - tag:GetResources
  - firehose:TagDeliveryStream
  - s3:GetBucketPolicy
  - s3:PutBucketPolicy
  Resource: "*"
- Effect: Allow
  Action: iam:CreateServiceLinkedRole
  Resource: arn:aws:iam::*:role/aws-service-role/vpc-lattice.amazonaws.com/AWSServiceRoleForVpcLattice
  Condition:
    StringLike:
      iam:AWSServiceName: vpc-lattice.amazonaws.com
- Effect: Allow
  Action: iam:CreateServiceLinkedRole
  Resource: arn:aws:iam::*:role/aws-service-role/delivery.logs.amazonaws.com/AWSServiceRoleForLogDelivery
  Condition:
    StringLike:
      iam:AWSServiceName: delivery.logs.amazonaws.com
`

// https://github.com/aws/aws-application-networking-k8s/blob/main/helm/values.yaml
const valuesTemplate = `---
fullnameOverride: gateway-api-controller
image:
  tag: {{ .Version }}
deployment:
  replicas: {{ .Replicas }}
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
{{- if .DefaultServiceNetwork }}
defaultServiceNetwork: {{ .DefaultServiceNetwork }}
{{- end }}
`
