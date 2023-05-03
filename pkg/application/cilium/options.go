package cilium

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type CiliumOptions struct {
	application.ApplicationOptions

	Overlay   bool
	Wireguard bool
}

func newOptions() (options *CiliumOptions, flags cmd.Flags) {
	options = &CiliumOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace: "kube-system",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.12.6",
				Latest:        "v1.12.6",
				PreviousChart: "1.11.6",
				Previous:      "v1.11.6",
			},
			// Cilium has many ServiceAccounts: cilium, etcd, operator, preflight, relay, ui
			DisableServiceAccountFlag: true,
		},
		Wireguard: false,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "overlay",
				Description: "run in overlay mode (remove VPC CNI first)",
			},
			Option: &options.Overlay,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "wireguard",
				Description: "enable wireguard transparent encryption",
			},
			Option: &options.Wireguard,
		},
	}

	return
}
