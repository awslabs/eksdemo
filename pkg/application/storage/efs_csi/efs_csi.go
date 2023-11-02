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
// Version: Latest is Chart 2.5.0, App v1.7.0 (as of 10/11/23)

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
					Name: "efs-csi-controller-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"},
			}),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name:           "efs-csi-node-irsa",
					ServiceAccount: "efs-csi-node-sa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/service-role/AmazonEFSCSIDriverPolicy"},
			}),
		},

		Options: &application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "efs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.5.0",
				Latest:        "v1.7.0",
				PreviousChart: "2.2.7",
				Previous:      "v1.4.0",
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
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
node:
  serviceAccount:
    name: efs-csi-node-sa
    annotations:
      {{ .IrsaAnnotationFor "efs-csi-node-sa" }}
`
