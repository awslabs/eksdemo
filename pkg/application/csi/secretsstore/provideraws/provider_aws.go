package provideraws

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/secretsmanager/latest/userguide/integrating_csi_driver.html
// GitHub:  https://github.com/aws/secrets-store-csi-driver-provider-aws
// Helm:    https://github.com/aws/secrets-store-csi-driver-provider-aws/tree/main/charts/secrets-store-csi-driver-provider-aws
// Repo:    https://gallery.ecr.aws/aws-secrets-manager/secrets-store-csi-driver-provider-aws
// Version: Latest is Chart 0.3.9 (as of 8/7/24)

func NewApp() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Parent:      "secrets",
			Name:        "store-csi-driver-provider-aws",
			Description: "AWS Secrets Manager and Config Provider for Secret Store CSI Driver",
			Aliases:     []string{"ascp", "provider-aws", "aws"},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "secrets-store-csi-driver-provider-aws",
			ReleaseName:   "secrets-store-csi-driver-provider-aws",
			RepositoryURL: "https://aws.github.io/secrets-store-csi-driver-provider-aws",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.3.9",
				Latest:        "1.0.r2-72-gfb78a36-2024.05.29.23.03",
				PreviousChart: "0.3.9",
				Previous:      "1.0.r2-72-gfb78a36-2024.05.29.23.03",
			},
			Namespace:      "kube-system",
			ServiceAccount: "secrets-store-csi-driver-provider-aws",
		},
	}
}

// https://github.com/aws/secrets-store-csi-driver-provider-aws/blob/main/charts/secrets-store-csi-driver-provider-aws/values.yaml
const valuesTemplate = `---
image:
  tag: {{ .Version }}
rbac:
  serviceAccountName: {{ .ServiceAccount }}
`
