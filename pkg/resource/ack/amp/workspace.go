package amp

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type WorkspaceOptions struct {
	resource.CommonOptions
	Alias string
}

func (o *WorkspaceOptions) SetName(name string) {
	o.Alias = name
}

func NewWorkspaceResource() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "amp-workspace",
			Description: "AMP Workspace",
			Aliases:     []string{"amp"},
			Args:        []string{"ALIAS"},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: workspaceYamlTemplate,
			},
		},

		Options: &WorkspaceOptions{
			CommonOptions: resource.CommonOptions{
				Name:          "ack-amp-workspace",
				Namespace:     "default",
				NamespaceFlag: true,
			},
		},
	}
}

const workspaceYamlTemplate = `---
apiVersion: prometheusservice.services.k8s.aws/v1alpha1
kind: Workspace
metadata:
  namespace: {{ .Namespace }}
  name: {{ .Alias }}
spec:
  alias: {{ .Alias }}
`
