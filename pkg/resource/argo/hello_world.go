package argo

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewHelloWorld() *resource.Resource {
	options := &resource.CommonOptions{
		Name:          "",
		Namespace:     "default",
		NamespaceFlag: true,
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "hello-world",
			Description: "Argo Hello World Workflow",
			Aliases:     []string{"helloworld", "hello"},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: helloWorldTemplate,
			},
		},

		Options: options,
	}
}

const helloWorldTemplate = `---
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: hello-world-
  namespace: {{ .Namespace }}
  labels:
    workflows.argoproj.io/archive-strategy: "false"
  annotations:
    workflows.argoproj.io/description: |
      This is a simple hello world example.
spec:
  entrypoint: whalesay
  templates:
  - name: whalesay
    container:
      image: docker/whalesay:latest
      command: [cowsay]
      args: ["hello world"]
`
