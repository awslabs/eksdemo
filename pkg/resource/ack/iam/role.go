package iam

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type RoleOptions struct {
	resource.CommonOptions
	Description string
}

func NewRoleResource() *resource.Resource {
	options := &RoleOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-iam-role",
			Namespace:     "default",
			NamespaceFlag: true,
		},
		Description: "ACK test role",
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "iam-role",
			Description: "IAM Role",
			Aliases:     []string{"role"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "description",
					Description: "role description",
				},
				Option: &options.Description,
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
apiVersion: iam.services.k8s.aws/v1alpha1
kind: Role
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  name: {{ .Name }}  
  description: {{ .Description }}
  assumeRolePolicyDocument: >
    {
      "Version":"2012-10-17",
      "Statement": [{
        "Effect":"Allow",
        "Principal": {
          "Service": [
            "ec2.amazonaws.com"
          ]
        },
        "Action": ["sts:AssumeRole"]
      }]
    }
`
