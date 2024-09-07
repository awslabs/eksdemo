package ec2

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type VpcOptions struct {
	resource.CommonOptions
	CidrBlocks []string
}

func NewVpcResource() *resource.Resource {
	options := &VpcOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-ec2-vpc",
			Namespace:     "default",
			NamespaceFlag: true,
		},
		CidrBlocks: []string{"10.0.0.0/16"},
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc",
			Description: "Virtual Private Cloud (VPC)",
			Aliases:     []string{"vpcs"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringSliceFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "cidrs",
					Description: "list of IPv4 CIDR blocks for your VPC",
				},
				Option: &options.CidrBlocks,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: vpcYamlTemplate,
			},
		},

		Options: options,
	}
}

const vpcYamlTemplate = `---
apiVersion: ec2.services.k8s.aws/v1alpha1
kind: VPC
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  cidrBlocks:
{{- range .CidrBlocks }}
  - {{ . }}
{{- end }}
  enableDNSHostnames: true
  enableDNSSupport: true
  tags:
  - key: Name
    value: {{ .Name }}
`
