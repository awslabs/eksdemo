package ascp

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

//

func NewApp() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "ascp",
			Description: "Example for AWS Secrets Manager and Config Provider for Secret Store CSI Driver",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "example-ascp-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-ascp",
			ResourceTemplate: &template.TextTemplate{
				Template: secretsProviderClassTemplate + serviceAccountTemplate + serviceAndDeploymentTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "ascp",
			ServiceAccount:            "nginx-deployment-sa",
		},
	}
}

// https://github.com/aws/secrets-store-csi-driver-provider-aws#usage
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - secretsmanager:GetSecretValue
  - secretsmanager:DescribeSecret
  Resource: arn:{{ .Partition }}:secretsmanager:{{ .Region }}:{{ .Account }}:secret:MySecret-??????
`
