package apigatewayv2_controller

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
// GitHub:  https://github.com/aws-controllers-k8s/apigatewayv2-controller
// Helm:    https://github.com/aws-controllers-k8s/apigatewayv2-controller/tree/main/helm
// Chart:   https://gallery.ecr.aws/aws-controllers-k8s/apigatewayv2-chart
// Repo:    https://gallery.ecr.aws/aws-controllers-k8s/apigatewayv2-controller
// Version: Latest is v1.0.16 (as of 9/5/24)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ack",
			Name:        "apigatewayv2-controller",
			Description: "ACK API Gateway v2 Controller",
			Aliases:     []string{"apigatewayv2", "apigwv2"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ack-apigatewayv2-controller-irsa",
				},
				// https://github.com/aws-controllers-k8s/apigatewayv2-controller/blob/main/config/iam/recommended-policy-arn
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonAPIGatewayAdministrator",
					"arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "ack-system",
			ServiceAccount: "ack-apigatewayv2-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.16",
				Latest:        "1.0.16",
				PreviousChart: "1.0.3",
				Previous:      "1.0.3",
			},
		},

		Installer: &installer.HelmInstaller{
			ReleaseName:   "ack-apigatewayv2-controller",
			RepositoryURL: "oci://public.ecr.aws/aws-controllers-k8s/apigatewayv2-chart",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

// https://github.com/aws-controllers-k8s/apigatewayv2-controller/blob/main/helm/values.yaml
const valuesTemplate = `---
image:
  tag: {{ .Version }}
fullnameOverride: ack-apigatewayv2-controller
aws:
  region: {{ .Region }}
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`
