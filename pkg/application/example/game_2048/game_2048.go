package game_2048

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// GitHub:   https://github.com/alexwhen/docker-2048
// Manifest: https://github.com/kubernetes-sigs/aws-load-balancer-controller/blob/main/docs/examples/2048/2048_full.yaml
// Repo:     https://gallery.ecr.aws/l6m2t8p7/docker-2048

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "game-2048",
			Description: "Example Game 2048",
			Aliases:     []string{"game2048", "2048"},
		},

		Installer: &installer.ManifestInstaller{
			AppName: "example-game-2048",
			ResourceTemplate: &template.TextTemplate{
				Template: gameManifestTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app
}
