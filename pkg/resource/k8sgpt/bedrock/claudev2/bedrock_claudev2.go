package claudev2

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type Options struct {
	resource.CommonOptions
	K8sGPTVersion string
}

func NewResource() *resource.Resource {
	options := &Options{
		CommonOptions: resource.CommonOptions{
			Name:          "k8sgpt-bedrock-claudev2",
			Namespace:     "k8sgpt",
			NamespaceFlag: true,
		},
		K8sGPTVersion: "v0.3.40",
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "bedrock-claudev2",
			Description: "Anthropic Claude v2 in Amazon Bedrock",
			Aliases:     []string{"claudev2"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "version",
					Description: "version of K8sGPT",
				},
				Option: &options.K8sGPTVersion,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},

		Options: options,
	}
}

const yamlTemplate = `---
apiVersion: core.k8sgpt.ai/v1alpha1
kind: K8sGPT
metadata:
  name: bedrock-claude-v2
  namespace: {{ .Namespace }}
spec:
  ai:
    enabled: true
    model: anthropic.claude-v2
    region: {{ .Region }}
    backend: amazonbedrock
  noCache: false
  repository: ghcr.io/k8sgpt-ai/k8sgpt
  version: {{ .K8sGPTVersion }}
`
