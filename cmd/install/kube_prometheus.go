package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/kube_prometheus/karpenter_dashboards"
	"github.com/awslabs/eksdemo/pkg/application/kube_prometheus/kube_prometheus_stack"
	"github.com/awslabs/eksdemo/pkg/application/kube_prometheus/kube_prometheus_stack_amp"
	"github.com/spf13/cobra"
)

var kubePrometheusApps []func() *application.Application

func NewInstallKubePrometheusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kube-prometheus",
		Short:   "End-to-end Cluster Monitoring with Prometheus",
		Aliases: []string{"kprom"},
	}

	// Don't show flag errors for `install kube-prometheus` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range kubePrometheusApps {
		cmd.AddCommand(a().NewInstallCmd())
	}

	return cmd
}

func NewUninstallKubePrometheusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kube-prometheus",
		Short:   "End-to-end Cluster Monitoring with Prometheus",
		Aliases: []string{"kprom"},
	}

	// Don't show flag errors for `uninstall kube-prometheus` without a subcommand
	cmd.DisableFlagParsing = true

	for _, a := range kubePrometheusApps {
		cmd.AddCommand(a().NewUninstallCmd())
	}

	return cmd
}

func init() {
	kubePrometheusApps = []func() *application.Application{
		karpenter_dashboards.NewApp,
		kube_prometheus_stack.NewApp,
		kube_prometheus_stack_amp.NewApp,
	}
}
