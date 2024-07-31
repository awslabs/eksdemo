package linkerdControlPlane

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AppOptions struct {
	application.ApplicationOptions
	trustAnchor string
	issuerCert string
	issuerKey string
}

func newOptions() (options *AppOptions, flags cmd.Flags) {
	options = &AppOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2024.7.3",
				PreviousChart: "2024.7.2",
			},
			Namespace: "linkerd",
		},
//		trustAnchor: "./pkg/application/linkerd/linkerd_control_plane/ca.crt",
//		issuerCert: "./pkg/application/linkerd/linkerd_control_plane/issuer.crt",
//		issuerKey: "./pkg/application/linkerd/linkerd_control_plane/issuer.key",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "trust-anchor",
				Description: "Path to TLS Certificate to use as the Trust Anchor",
			},
			Option: &options.trustAnchor,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "issuer-cert",
				Description: "Path to TLS Certificate to use as the Issuer",
			},
			Option: &options.issuerCert,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "issuer-key",
				Description: "Path to TLS Key to use as the Issuer",
			},
			Option: &options.issuerKey,
		},
	}
	return
}
