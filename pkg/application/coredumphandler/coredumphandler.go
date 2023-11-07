package coredumphandler

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/s3_bucket"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/IBM/core-dump-handler
// Helm:    https://ibm.github.io/core-dump-handler/
// Values:  https://github.com/IBM/core-dump-handler/blob/main/charts/core-dump-handler/values.aws.sts.yaml
// Repo:    https://quay.io/repository/icdh/core-dump-handler
// Version: Latest is Chart/App v8.10.0 (as of 10/24/23)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Name:        "core-dump-handler",
			Description: "Automatically saves core dumps to S3",
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "core-dump-handler-irsa",
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
			ChartName:     "core-dump-handler",
			ReleaseName:   "core-dump-handler",
			RepositoryURL: "https://ibm.github.io/core-dump-handler/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}

	return app
}

const valuesTemplate = `---
# https://github.com/IBM/core-dump-handler/blob/main/charts/core-dump-handler/values.aws.sts.yaml
daemonset:
  includeCrioExe: {{ not .IncludeCrioExe }}
  vendor: rhel7 # EKS EC2 images have an old libc=2.26
  s3BucketName: eksdemo-{{ .Account }}-coredumphandler
  s3Region: {{ .Region }}

serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
`

// https://github.com/IBM/core-dump-handler/blob/main/charts/core-dump-handler/README.md#eks-setup-for-gitops-pipelines-eksctl-or-similar
const policyDocument = `
Version: '2012-10-17'
Statement:
- Effect: Allow
  Action:
  - s3:*
  Resource: [
   "arn:{{ .Partition }}:s3:::eksdemo-{{ .Account }}-coredumphandler",
   "arn:{{ .Partition }}:s3:::eksdemo-{{ .Account }}-coredumphandler/*"
  ]
`
