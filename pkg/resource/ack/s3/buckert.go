package s3

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewBucketResource() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "s3-bucket",
			Description: "Amazon S3 bucket",
			Aliases:     []string{"s3"},
			CreateArgs:  []string{"NAME"},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: bucketYamlTemplate,
			},
		},

		Options: &resource.CommonOptions{
			Name:          "ack-s3-bucket",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}
}

const bucketYamlTemplate = `---
apiVersion: s3.services.k8s.aws/v1alpha1
kind: Bucket
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  name: {{ .Name }}
`
