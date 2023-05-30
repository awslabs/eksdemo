package ebs_csi

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/kubernetes-sigs/aws-ebs-csi-driver/tree/master/docs
// GitHub:  https://github.com/kubernetes-sigs/aws-ebs-csi-driver
// Helm:    https://github.com/kubernetes-sigs/aws-ebs-csi-driver/tree/master/charts/aws-ebs-csi-driver
// Repo:    gallery.ecr.aws/ebs-csi-driver/aws-ebs-csi-driver
// Version: Latest is v1.19.0, Chart 2.19.0 (as of 5/28/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "storage",
			Name:        "ebs-csi",
			Description: "Amazon EBS CSI driver",
			Aliases:     []string{"ebscsi", "ebs"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ebs-csi-irsa",
				},
				PolicyType: irsa.PolicyARNs,
				Policy:     []string{"arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"},
			}),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "aws-ebs-csi-driver",
			ReleaseName:   "storage-ebs-csi",
			RepositoryURL: "https://kubernetes-sigs.github.io/aws-ebs-csi-driver",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
image:
  tag: {{ .Version }}
controller:
  region: {{ .Region }}
  replicaCount: {{ .Replicas }}
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
{{- if not .NoStorageClasses }}
storageClasses:
- name: gp3
  allowVolumeExpansion: true
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  parameters:
    csi.storage.k8s.io/fstype: ext4
    type: gp3
  volumeBindingMode: WaitForFirstConsumer
- name: gp3-encrypted
  allowVolumeExpansion: true
  parameters:
    csi.storage.k8s.io/fstype: ext4
    encrypted: "true"
    type: gp3
  volumeBindingMode: WaitForFirstConsumer
{{- end }}
`
