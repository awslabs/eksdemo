package vpc_lattice_controller

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
// Version: Latest is v0.0.16 (as of 9/13/23)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
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

		Installer: &installer.HelmInstaller{
			ReleaseName:   "vpc-lattice-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-application-networking-k8s/aws-gateway-controller-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		PostInstallResources: []*resource.Resource{
			gatewayClass(),
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

// https://github.com/aws/aws-application-networking-k8s/blob/main/config/iam/recommended-inline-policy.json
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - vpc-lattice:*
  - iam:CreateServiceLinkedRole
  - ec2:DescribeVpcs
  - ec2:DescribeSubnets
  - ec2:DescribeTags
  Resource: "*"
`

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
`
