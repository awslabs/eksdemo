package efs

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type FileSystemOptions struct {
	resource.CommonOptions
	ThroughputMode string
}

func NewFileSystemResource() *resource.Resource {
	options := &FileSystemOptions{
		CommonOptions: resource.CommonOptions{
			Namespace:     "default",
			NamespaceFlag: true,
		},
		ThroughputMode: "bursting",
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "efs-filesystem",
			Description: "EFS File System",
			Aliases:     []string{"efs"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "throughput-mode",
					Description: "specifies the throughput mode for the file system",
				},
				Option:  &options.ThroughputMode,
				Choices: []string{"elastic", "provisioned", "bursting"},
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: fileSystemYamlTemplate,
			},
		},

		Options: options,
	}
}

const fileSystemYamlTemplate = `---
apiVersion: efs.services.k8s.aws/v1alpha1
kind: FileSystem
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  encrypted: true
  fileSystemProtection:
    replicationOverwriteProtection: ENABLED
  performanceMode: generalPurpose
  throughputMode: {{ .ThroughputMode }}
  tags:
  - key: Name
    value: {{ .Name }}
`
