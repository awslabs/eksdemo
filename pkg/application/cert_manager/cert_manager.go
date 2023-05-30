package cert_manager

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://cert-manager.io/docs/
// GitHub:  https://github.com/cert-manager/cert-manager
// Helm:    https://github.com/cert-manager/cert-manager/tree/master/deploy/charts/cert-manager
// Repo:    quay.io/jetstack/cert-manager-controller
// Version: Latest is v1.12.1 (as of 5/30/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "cert-manager",
			Description: "Cloud Native Certificate Management",
			Aliases:     []string{"certmanager"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "cert-manager-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "cert-manager",
			ServiceAccount: "cert-manager",
			DefaultVersion: &application.LatestPrevious{
				// For Helm Chart version: https://artifacthub.io/packages/helm/cert-manager/cert-manager
				LatestChart:   "1.12.1",
				Latest:        "v1.12.1",
				PreviousChart: "1.11.0",
				Previous:      "v1.11.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cert-manager",
			ReleaseName:   "cert-manager",
			RepositoryURL: "https://charts.jetstack.io",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		PostInstallResources: []*resource.Resource{
			clusterIssuer(),
		},
	}
	return app
}

const valuesTemplate = `---
installCRDs: true
replicaCount: 1
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
image:
  tag: {{ .Version }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - route53:GetChange
  Resource: arn:aws:route53:::change/*
- Effect: Allow
  Action:
  - route53:ChangeResourceRecordSets
  - route53:ListResourceRecordSets
  Resource: arn:aws:route53:::hostedzone/*
- Effect: Allow
  Action: route53:ListHostedZonesByName
  Resource: "*"
`
