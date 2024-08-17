package k8sgpt

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.k8sgpt.ai/getting-started/in-cluster-operator/
// GitHub:  https://github.com/k8sgpt-ai/k8sgpt-operator
// Helm:    https://github.com/k8sgpt-ai/k8sgpt-operator/tree/main/chart/operator
// Repo:    ghcr.io/k8sgpt-ai/k8sgpt-operator
// Version: Latest is v0.1.7 (as of 8/17/24)

func NewApp() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Name:        "k8sgpt-operator",
			Description: "K8sGPT Operator",
			Aliases:     []string{"k8sgpt"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "k8sgpt-operator-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "k8sgpt-operator",
			ReleaseName:   "k8sgpt-operator",
			RepositoryURL: "https://charts.k8sgpt.ai",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: &application.ApplicationOptions{
			Namespace: "k8sgpt",
			// This isn't the SA of the operator, it's the SA of the k8sgpt instance created by the operator
			// This is for IRSA to give permission to the k8sgpt instance to access Bedrock APIs
			ServiceAccount: "k8sgpt",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "0.1.7",
				Latest:        "v0.1.7",
				PreviousChart: "0.1.6",
				Previous:      "v0.1.6",
			},
		},
	}
}

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - bedrock:InvokeModel
  Resource: "*"
`

// https://github.com/k8sgpt-ai/k8sgpt-operator/blob/main/chart/operator/values.yaml
const valuesTemplate = `---
serviceAccount:
  name: {{ .ServiceAccount }}
  # -- Annotations for the managed k8sgpt workload service account
  annotations:
    {{ .IrsaAnnotation }}
controllerManager:
  manager:
    image:
      tag: {{ .Version }}
`
