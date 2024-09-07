package ec2

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type SecurityGroupOptions struct {
	resource.CommonOptions
	Description string
}

func NewSecurityGroupResource() *resource.Resource {
	options := &SecurityGroupOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-ec2-security-group",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "security-group",
			Description: "EC2 Security Group",
			Aliases:     []string{"security-groups", "sg"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "desc",
					Description: "description",
					Required:    true,
				},
				Option: &options.Description,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: securityGroupYamlTemplate,
			},
		},

		Options: options,
	}
}

const securityGroupYamlTemplate = `---
apiVersion: ec2.services.k8s.aws/v1alpha1
kind: SecurityGroup
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  description: {{ .Description }}
  name: {{ .Name }}
  tagSpecifications:
  - resourceType: security-group
    tags:
    - key: Name
      value: {{ .Name }}
  vpcID: {{ .Cluster.ResourcesVpcConfig.VpcId }}
`
