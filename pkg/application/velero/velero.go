package velero

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/s3_bucket"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://velero.io/docs/
// GitHub:  https://github.com/vmware-tanzu/velero
// GitHub:  https://github.com/vmware-tanzu/velero-plugin-for-aws
// Helm:    https://github.com/vmware-tanzu/helm-charts/tree/main/charts/velero
// Repo:    velero/velero
// Version: Latest is Chart 5.0.2, app v1.11.1 (as of 9/13/23)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "velero",
			Description: "Backup and Migrate Kubernetes Applications",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "velero-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: policyDocument,
				},
			}),
			s3_bucket.NewResourceWithOptions(options.BucketOptions),
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "velero",
			ReleaseName:   "velero",
			RepositoryURL: "https://vmware-tanzu.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}
}

// https://github.com/vmware-tanzu/velero-plugin-for-aws#set-permissions-for-velero
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - ec2:DescribeVolumes
  - ec2:DescribeSnapshots
  - ec2:CreateTags
  - ec2:CreateVolume
  - ec2:CreateSnapshot
  - ec2:DeleteSnapshot
  Resource: "*"
- Effect: Allow
  Action:
  - s3:GetObject
  - s3:DeleteObject
  - s3:PutObject
  - s3:AbortMultipartUpload
  - s3:ListMultipartUploadParts
  Resource: arn:aws:s3:::eksdemo-{{ .Account }}-velero/*
- Effect: Allow
  Action: s3:ListBucket
  Resource: arn:aws:s3:::eksdemo-{{ .Account }}-velero
`

const valuesTemplate = `---
image:
  tag: {{ .Version }}
initContainers:
- name: velero-plugin-for-aws
  image: velero/velero-plugin-for-aws:{{ .PluginVersion }}
  volumeMounts:
  - mountPath: /target
    name: plugins
configuration:
  backupStorageLocation:
  - provider: aws
    bucket: eksdemo-{{ .Account }}-velero
  volumeSnapshotLocation:
  - provider: aws
    config:
      region: {{ .Region }}
serviceAccount:
  server:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
credentials:
  useSecret: false
`
