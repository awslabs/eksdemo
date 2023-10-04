package kube_prometheus_stack_amp

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
)

const AmpAliasSuffix = "kube-prometheus"

type KubePrometheusStackAmpOptions struct {
	application.ApplicationOptions
	*amp_workspace.AmpWorkspaceOptions

	AmpEndpoint           string
	DisableGrafana        bool
	GrafanaAdminPassword  string
	GrafanaServiceAccount string
}

func NewOptions() (options *KubePrometheusStackAmpOptions, flags cmd.Flags) {
	options = &KubePrometheusStackAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "51.2.0",
				Latest:        "v0.68.0",
				PreviousChart: "46.6.0",
				Previous:      "v0.65.1",
			},
			DisableServiceAccountFlag:    true,
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "monitoring",
			ServiceAccount:               "prometheus-prometheus",
		},
		GrafanaServiceAccount: "prometheus-grafana",
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "grafana-pass",
				Description: "grafana admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.GrafanaAdminPassword,
		},
	}
	return
}

func (o *KubePrometheusStackAmpOptions) PreDependencies(application.Action) error {
	o.AmpWorkspaceOptions.Alias = fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix)
	return nil
}

func (o *KubePrometheusStackAmpOptions) PreInstall() error {
	if o.DryRun {
		o.AmpEndpoint = "<-amp_endpoint_url_will_go_here->"
		return nil
	}
	ampGetter := amp_workspace.NewGetter(aws.NewAMPClient())

	workspace, err := ampGetter.GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix))
	if err != nil {
		return fmt.Errorf("failed to lookup AMP endpoint to use in Helm chart for remoteWrite url: %w", err)
	}

	o.AmpEndpoint = *workspace.Workspace.PrometheusEndpoint

	return nil
}
