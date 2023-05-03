package addon

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "addon",
			Description: "EKS Managed Addon",
			Aliases:     []string{"addons"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "addon",
			ConfigTemplate: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlTemplate,
			},
			DeleteFlags: &template.TextTemplate{
				Template: deleteFlagsTemplate,
			},
		},
	}

	res.Options, res.CreateFlags = NewOptions()

	return res
}

const eksctlTemplate = `
addons:
- name: {{ .Name }}
{{- if .Version }}
  version: {{ .Version }}
{{- end }}
`
const deleteFlagsTemplate = `--name {{ .Name }} --cluster {{ .ClusterName }} --region {{ .Region }} --preserve`
