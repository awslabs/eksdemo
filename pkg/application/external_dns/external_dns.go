package external_dns

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md
// GitHub:  https://github.com/kubernetes-sigs/external-dns
// Helm:    https://github.com/kubernetes-sigs/external-dns/tree/master/charts/external-dns
// Repo:    registry.k8s.io/external-dns/external-dns
// Version: Latest is Chart 1.13.1, App v0.13.6 (as of 10/2/23)

func New() *application.Application {
	return &application.Application{
		Command: cmd.Command{
			Name:        "external-dns",
			Description: "ExternalDNS",
			Aliases:     []string{"externaldns", "edns"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "external-dns-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "external-dns",
			ServiceAccount: "external-dns",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.13.1",
				Latest:        "v0.13.6",
				PreviousChart: "1.12.2",
				Previous:      "v0.13.4",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "external-dns",
			ReleaseName:   "external-dns",
			RepositoryURL: "https://kubernetes-sigs.github.io/external-dns",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
}

// https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#iam-policy
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - route53:ChangeResourceRecordSets
  Resource:
  - arn:{{ .Partition }}:route53:::hostedzone/*
- Effect: Allow
  Action:
  - route53:ListHostedZones
  - route53:ListResourceRecordSets
  - route53:ListTagsForResource
  Resource:
  - "*"
`

const valuesTemplate = `---
image:
  tag: {{ .Version }}
provider: aws
registry: txt
serviceAccount:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount }}
txtOwnerId: {{ .ClusterName }}
`
