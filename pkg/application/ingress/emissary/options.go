package emissary

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
)

type EmissaryOptions struct {
	application.ApplicationOptions

	Replicas int
}

const crdUrl = "https://app.getambassador.io/yaml/emissary/%s/emissary-crds.yaml"

func newOptions() (options *EmissaryOptions, flags cmd.Flags) {
	options = &EmissaryOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "emissary",
			ServiceAccount: "emissary-ingress",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "8.0.0",
				Latest:        "3.0.0",
				PreviousChart: "8.0.0",
				Previous:      "3.0.0",
			},
		},
		Replicas: 1,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of Ambassador replicas",
			},
			Option: &options.Replicas,
		},
	}

	return
}

func (o *EmissaryOptions) PreInstall() error {
	url := fmt.Sprintf(crdUrl, o.Version)

	if o.DryRun {
		fmt.Println("Preinstall will install CRDs from: " + url)
		return nil
	}

	fmt.Println("Installing CRDs from: " + url)

	return kubernetes.CreateResources(o.KubeContext(), url)
}
