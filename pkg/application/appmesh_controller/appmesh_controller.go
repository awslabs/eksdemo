package appmesh_controller

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/app-mesh/latest/userguide/
// Docs:    https://aws.github.io/aws-app-mesh-controller-for-k8s/
// GitHub:  https://github.com/aws/aws-app-mesh-controller-for-k8s
// Helm:    https://github.com/aws/eks-charts/tree/master/stable/appmesh-controller
// Repo:    602401143452.dkr.ecr.us-west-2.amazonaws.com/amazon/appmesh-controller
// Version: Latest is v1.7.0 (as of 11/04/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "appmesh-controller",
			Description: "AWS App Mesh Controller",
			Aliases:     []string{"appmesh", "app-mesh"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "appmesh-controller-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "appmesh-system",
			ServiceAccount: "appmesh-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.7.0",
				Latest:        "v1.7.0",
				PreviousChart: "1.5.0",
				Previous:      "v1.5.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "appmesh-controller",
			ReleaseName:   "appmesh-controller",
			RepositoryURL: "https://aws.github.io/eks-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
region: {{ .Region }}
accountId: {{ .Account }}
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
image:
  tag: {{ .Version }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - appmesh:ListVirtualRouters
  - appmesh:ListVirtualServices
  - appmesh:ListRoutes
  - appmesh:ListGatewayRoutes
  - appmesh:ListMeshes
  - appmesh:ListVirtualNodes
  - appmesh:ListVirtualGateways
  - appmesh:DescribeMesh
  - appmesh:DescribeVirtualRouter
  - appmesh:DescribeRoute
  - appmesh:DescribeVirtualNode
  - appmesh:DescribeVirtualGateway
  - appmesh:DescribeGatewayRoute
  - appmesh:DescribeVirtualService
  - appmesh:CreateMesh
  - appmesh:CreateVirtualRouter
  - appmesh:CreateVirtualGateway
  - appmesh:CreateVirtualService
  - appmesh:CreateGatewayRoute
  - appmesh:CreateRoute
  - appmesh:CreateVirtualNode
  - appmesh:UpdateMesh
  - appmesh:UpdateRoute
  - appmesh:UpdateVirtualGateway
  - appmesh:UpdateVirtualRouter
  - appmesh:UpdateGatewayRoute
  - appmesh:UpdateVirtualService
  - appmesh:UpdateVirtualNode
  - appmesh:DeleteMesh
  - appmesh:DeleteRoute
  - appmesh:DeleteVirtualRouter
  - appmesh:DeleteGatewayRoute
  - appmesh:DeleteVirtualService
  - appmesh:DeleteVirtualNode
  - appmesh:DeleteVirtualGateway
  Resource: "*"
- Effect: Allow
  Action:
  - iam:CreateServiceLinkedRole
  Resource: arn:aws:iam::*:role/aws-service-role/appmesh.amazonaws.com/AWSServiceRoleForAppMesh
  Condition:
    StringLike:
      iam:AWSServiceName:
      - appmesh.amazonaws.com
- Effect: Allow
  Action:
  - acm:ListCertificates
  - acm:DescribeCertificate
  - acm-pca:DescribeCertificateAuthority
  - acm-pca:ListCertificateAuthorities
  Resource: "*"
- Effect: Allow
  Action:
  - servicediscovery:CreateService
  - servicediscovery:DeleteService
  - servicediscovery:GetService
  - servicediscovery:GetInstance
  - servicediscovery:RegisterInstance
  - servicediscovery:DeregisterInstance
  - servicediscovery:ListInstances
  - servicediscovery:ListNamespaces
  - servicediscovery:ListServices
  - servicediscovery:GetInstancesHealthStatus
  - servicediscovery:UpdateInstanceCustomHealthStatus
  - servicediscovery:GetOperation
  - route53:GetHealthCheck
  - route53:CreateHealthCheck
  - route53:UpdateHealthCheck
  - route53:ChangeResourceRecordSets
  - route53:DeleteHealthCheck
  Resource: "*"
`
