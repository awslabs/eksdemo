package efs_csi

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/aws-efs-csi-driver/tree/master/docs
// GitHub:  https://github.com/kubernetes-sigs/aws-efs-csi-driver
// Helm:    https://github.com/kubernetes-sigs/aws-efs-csi-driver/tree/master/charts/aws-efs-csi-driver
// Repo:    amazon/aws-efs-csi-driver
// Version: Latest is Chart 2.2.7, App v1.4.0 (as of 07/31/22)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "storage",
			Name:        "efs-csi",
			Description: "Amazon EFS CSI driver",
			Aliases:     []string{"efscsi", "efs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "efs-csi-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "efs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.2.7",
				Latest:        "v1.4.0",
				PreviousChart: "2.2.6",
				Previous:      "v1.3.8",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-efs-csi-driver",
			ReleaseName:   "storage-efs-csi",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-efs-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	return app
}

const valuesTemplate = `---
replicaCount: 1
image:
  tag: {{ .Version }}
controller:
  serviceAccount:
    annotations:
      {{ .IrsaAnnotation }}
    name: {{ .ServiceAccount }}
`

const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - elasticfilesystem:DescribeAccessPoints
  - elasticfilesystem:DescribeFileSystems
  - elasticfilesystem:DescribeMountTargets
  - ec2:DescribeAvailabilityZones
  Resource: "*"
- Effect: Allow
  Action:
  - elasticfilesystem:CreateAccessPoint
  Resource: "*"
  Condition:
    StringLike:
      aws:RequestTag/efs.csi.aws.com/cluster: 'true'
- Effect: Allow
  Action: elasticfilesystem:TagResource
  Resource: "*"
  Condition:
    StringLike:
	aws:ResourceTag/efs.csi.aws.com/cluster: 'true'
- Effect: Allow
  Action: elasticfilesystem:DeleteAccessPoint
  Resource: "*"
  Condition:
    StringEquals:
      aws:ResourceTag/efs.csi.aws.com/cluster: 'true'
`
