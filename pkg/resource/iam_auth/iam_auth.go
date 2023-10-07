package iam_auth

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type IamAuthOptions struct {
	resource.CommonOptions
	Groups   []string
	RoleName template.Template
	Username string
}

func NewResourceWithOptions(options *IamAuthOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-auth",
			Description: "IAM User or Role for Authentication",
		},

		Manager: &eksctl.ResourceManager{
			Resource: "iamidentitymapping",
			ConfigTemplate: &template.TextTemplate{
				Template: eksctl.EksctlHeader + eksctlTemplate,
			},
			DeleteCommand: &template.TextTemplate{
				Template: deleteCommandTemplate,
			},
		},
	}

	res.Options = options

	return res
}

const arn = `arn:{{ .Partition }}:iam::{{ .Account }}:role/{{ .RoleName.Render . }}`

const eksctlTemplate = `
iamIdentityMappings:
- arn: ` + arn + `
  groups:
{{- range .Groups }}
  - {{ . }}
{{- end }}
  username: {{ .Username }}
  noDuplicateARNs: true # prevents shadowing of ARNs
`
const deleteCommandTemplate = `--arn ` + arn + ` --cluster {{ .ClusterName }} --region {{ .Region }}`
