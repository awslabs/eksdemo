package eks_workshop

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// GitHub:   https://github.com/aws-containers/ecsdemo-nodejs
// GitHub:   https://github.com/aws-containers/ecsdemo-crystal
// GitHub:   https://github.com/aws-containers/ecsdemo-frontend
// Manifest: https://github.com/aws-containers/ecsdemo-nodejs/tree/main/kubernetes
// Manifest: https://github.com/aws-containers/ecsdemo-crystal/tree/main/kubernetes
// Manifest: https://github.com/aws-containers/ecsdemo-frontend/tree/main/kubernetes
// Repo:     https://gallery.ecr.aws/aws-containers/ecsdemo-nodejs
// Repo:     https://gallery.ecr.aws/aws-containers/ecsdemo-crystal
// Repo:     https://gallery.ecr.aws/aws-containers/ecsdemo-frontend

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "eks-workshop",
			Description: "EKS Workshop Example Microservices",
			Aliases:     []string{"eksworkshop"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-eks-workshop",
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
		},
	}

	app.Options, app.Flags = NewOptions()

	return app
}
