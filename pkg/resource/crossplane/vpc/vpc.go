package vpc

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type Options struct {
	resource.CommonOptions
	CidrBlock string
}

func NewResource() *resource.Resource {
	options := &Options{
		CommonOptions: resource.CommonOptions{
			Name:          "crossplane-vpc",
			Namespace:     "default",
			NamespaceFlag: true,
		},
		CidrBlock: "10.0.0.0/16",
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "cidr",
				Description: "IPv4 CIDR block for your VPC",
			},
			Option: &options.CidrBlock,
		},
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "vpc",
			Description: "Virtual Private Cloud",
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: flags,

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},

		Options: options,
	}
}

const yamlTemplate = `---
apiVersion: ec2.aws.upbound.io/v1beta1
kind: VPC
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  forProvider:
    region: {{ .Region }}
    cidrBlock: {{ .CidrBlock }}
    enableDnsSupport: true
    enableDnsHostnames: true
    tags:
      Name: {{ .Name }}
`
