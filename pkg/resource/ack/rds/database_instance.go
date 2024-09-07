package rds

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type DatabaseInstanceOptions struct {
	resource.CommonOptions
	Password string
	Storage  int
}

func NewDatabaseInstanceResource() *resource.Resource {
	options := &DatabaseInstanceOptions{
		CommonOptions: resource.CommonOptions{
			Namespace:     "default",
			NamespaceFlag: true,
		},
		Storage: 5,
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "rds-database-instance",
			Description: "RDS Database Instance",
			Aliases:     []string{"rds"},
			CreateArgs:  []string{"NAME"},
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "password",
					Description: "database instance root password",
					Required:    true,
					Shorthand:   "P",
				},
				Option: &options.Password,
			},
			&cmd.IntFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "storage",
					Description: "storage in gibibytes (GiB) to allocate for the DB instance",
				},
				Option: &options.Storage,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: databaseInstanceYamlTemplate,
			},
		},

		Options: options,
	}
}

const databaseInstanceYamlTemplate = `---
apiVersion: v1
kind: Secret
metadata:
  name: dbpass
  namespace: {{ .Namespace }}
data:
  password: {{ .Password | b64enc }}
---
apiVersion: rds.services.k8s.aws/v1alpha1
kind: DBInstance
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  allocatedStorage: {{ .Storage }}
  dbInstanceClass: db.t3.medium
  dbInstanceIdentifier: {{ .Name }}
  engine: mysql
  masterUsername: root
  masterUserPassword:
    name: dbpass
    key: password
`
