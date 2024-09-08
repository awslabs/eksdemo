package ecr

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type EcrOptions struct {
	resource.CommonOptions
	ImmutableTags bool
	ScanOnPush    bool
}

func NewRepositoryResource() *resource.Resource {
	options := &EcrOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-ecr-repo",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "ecr-repo",
			Description: "ECR Repository",
			Aliases:     []string{"ecr"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.BoolFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "immutable-tags",
					Description: "all image tags within the repository will be immutable",
				},
				Option: &options.ImmutableTags,
			},
			&cmd.BoolFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "scan-on-push",
					Description: "scan image when it's pushed to the repository",
				},
				Option: &options.ScanOnPush,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: repoYamlTemplate,
			},
		},

		Options: options,
	}
}

const repoYamlTemplate = `---
apiVersion: ecr.services.k8s.aws/v1alpha1
kind: Repository
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  name: {{ .Name }}
  imageScanningConfiguration:
    scanOnPush: {{ .ScanOnPush }}
{{- if .ImmutableTags }}
  imageTagMutability: IMMUTABLE
{{- end }}
`
