package amg_workspace

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	options, createFlags, deleteFlags := NewOptions()

	res := NewResourceWithOptions(options)
	res.CreateFlags = createFlags
	res.DeleteFlags = deleteFlags

	return res
}

func NewResourceWithOptions(options *AmgOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amg-workspace",
			Description: "Amazon Managed Grafana Workspace",
			Aliases:     []string{"amg"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{
			AssumeRolePolicyTemplate: template.TextTemplate{
				Template: assumeRolePolicyTemplate,
			},
		},
	}

	res.Options = options

	return res
}

const assumeRolePolicyTemplate = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "grafana.amazonaws.com"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "aws:SourceAccount": "{{ .Account }}"
        },
        "StringLike": {
          "aws:SourceArn": "arn:aws:grafana:{{ .Region }}:{{ .Account }}:/workspaces/*"
        }
      }
    }
  ]
}
`
const rolePolicName = `eksdemo-AMP-Policy`

const rolePolicyDoc = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "aps:ListWorkspaces",
        "aps:DescribeWorkspace",
        "aps:QueryMetrics",
        "aps:GetLabels",
        "aps:GetSeries",
        "aps:GetMetricMetadata"
      ],
      "Resource": "*"
    }
  ]
}
`
