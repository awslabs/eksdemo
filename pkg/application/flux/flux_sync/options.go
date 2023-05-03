package flux_sync

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type FluxSyncOptions struct {
	*application.ApplicationOptions
	GitUrl            string
	KustomizationPath string
	TargetNamespace   string
}

func NewOptions() (options *FluxSyncOptions, flags cmd.Flags) {
	options = &FluxSyncOptions{
		ApplicationOptions: &application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.0.0",
				PreviousChart: "0.4.2",
			},
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "default",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "git-url",
				Description: "git repository url to sync with",
				Required:    true,
			},
			Option: &options.GitUrl,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "kustomization-path",
				Description: "path to the directory containing the kustomization.yaml",
			},
			Option: &options.KustomizationPath,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-namespace",
				Description: "sets or overrides the namespace in kustomization.yaml",
			},
			Option: &options.TargetNamespace,
		},
	}
	return
}
