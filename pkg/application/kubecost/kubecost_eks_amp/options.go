package kubecost_eks_amp

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
)

const AmpAliasSuffix = "kubecost"

type KubecostEksAmpOptions struct {
	application.ApplicationOptions
	*amp_workspace.AmpWorkspaceOptions

	AmpEndpoint        string
	AmpId              string
	DisablePrometheus  bool
	EnableNodeExporter bool
}

func newOptions() (options *KubecostEksAmpOptions, flags cmd.Flags) {
	options = &KubecostEksAmpOptions{
		ApplicationOptions: application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kubecost",
			ServiceAccount:               "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.103.3",
				Latest:        "1.103.3",
				PreviousChart: "1.100.0",
				Previous:      "1.100.0",
			},
		},
		AmpWorkspaceOptions: &amp_workspace.AmpWorkspaceOptions{
			CommonOptions: resource.CommonOptions{
				Name: "kubecost-amazon-managed-prometheus",
			},
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-prometheus",
				Description: "don't install prometheus",
			},
			Option: &options.DisablePrometheus,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "node-exporter",
				Description: "install prometheus node exporter (not installed by default)",
			},
			Option: &options.EnableNodeExporter,
		},
	}
	return
}

func (o *KubecostEksAmpOptions) PreDependencies(application.Action) error {
	o.AmpWorkspaceOptions.Alias = fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix)
	return nil
}

func (o *KubecostEksAmpOptions) PreInstall() error {
	o.AmpEndpoint = "<-amp_endpoint_url_will_go_here->"
	o.AmpId = "<-amp_id_will_go_here->"

	workspace, err := amp_workspace.NewGetter(aws.NewAMPClient()).GetAmpByAlias(fmt.Sprintf("%s-%s", o.ClusterName, AmpAliasSuffix))
	if err != nil {
		if o.DryRun {
			return nil
		}
		return fmt.Errorf("failed to lookup AMP to use in Helm chart values file: %w", err)
	}

	o.AmpEndpoint = awssdk.ToString(workspace.Workspace.PrometheusEndpoint)
	o.AmpId = awssdk.ToString(workspace.Workspace.WorkspaceId)

	return nil
}
