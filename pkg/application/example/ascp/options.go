package ascp

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type Options struct {
	application.ApplicationOptions
	JSONFormat bool
	K8sSecret  bool
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "ascp",
			ServiceAccount:            "nginx-deployment-sa",
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "json-format",
				Description: "mount key/value pairs from a secret in json format",
			},
			Option: &options.JSONFormat,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "k8s-secret",
				Description: "create a Kubernetes Secret to mirror the mounted secret",
			},
			Option: &options.K8sSecret,
		},
	}

	return
}
