package fargate_profile

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	options, createFlags := NewOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "fargate-profile",
			Description: "EKS Fargate Profile",
			Aliases:     []string{"fargate-profiles", "fargateprofiles", "fargateprofile", "fargate", "fp"},
			CreateArgs:  []string{"NAME"},
			Args:        []string{"NAME"},
		},

		CreateFlags: createFlags,

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "fargateprofile",
			ConfigTemplate: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlTemplate,
			},
			DeleteCommand: &template.TextTemplate{
				Template: deleteCommandTemplate,
			},
		},

		Options: options,
	}
}

const eksctlTemplate = `
fargateProfiles:
- name: {{ .Name }}
  selectors:
{{- range .Namespaces }}
  - namespace: {{ . }}
{{- end }}
`
const deleteCommandTemplate = `--name {{ .Name }} --cluster {{ .ClusterName }} --region {{ .Region }}`
