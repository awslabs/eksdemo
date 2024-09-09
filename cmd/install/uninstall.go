package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/adot"
	"github.com/awslabs/eksdemo/pkg/application/argo/cd"
	"github.com/awslabs/eksdemo/pkg/application/autoscaling/cluster_autoscaler"
	awslbc "github.com/awslabs/eksdemo/pkg/application/aws/lbc"
	"github.com/awslabs/eksdemo/pkg/application/aws_fluent_bit"
	"github.com/awslabs/eksdemo/pkg/application/certmanager"
	"github.com/awslabs/eksdemo/pkg/application/cilium"
	"github.com/awslabs/eksdemo/pkg/application/consul"
	"github.com/awslabs/eksdemo/pkg/application/coredumphandler"
	"github.com/awslabs/eksdemo/pkg/application/externaldns"
	"github.com/awslabs/eksdemo/pkg/application/falco"
	"github.com/awslabs/eksdemo/pkg/application/goldilocks"
	"github.com/awslabs/eksdemo/pkg/application/harbor"
	"github.com/awslabs/eksdemo/pkg/application/headlamp"
	"github.com/awslabs/eksdemo/pkg/application/k8sgpt"
	"github.com/awslabs/eksdemo/pkg/application/karpenter"
	"github.com/awslabs/eksdemo/pkg/application/keycloak_amg"
	"github.com/awslabs/eksdemo/pkg/application/kube_state_metrics"
	"github.com/awslabs/eksdemo/pkg/application/metrics_server"
	"github.com/awslabs/eksdemo/pkg/application/prometheus_node_exporter"
	sparkoperator "github.com/awslabs/eksdemo/pkg/application/spark"
	"github.com/awslabs/eksdemo/pkg/application/storage/ebs_csi"
	"github.com/awslabs/eksdemo/pkg/application/vault"
	"github.com/awslabs/eksdemo/pkg/application/velero"
	"github.com/awslabs/eksdemo/pkg/application/vpclattice"
	"github.com/spf13/cobra"
)

func NewUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "uninstall",
		Short:   "Uninstall application and delete dependencies",
		Aliases: []string{"uninst"},
	}

	// Don't show flag errors for uninstall without a subcommand
	cmd.DisableFlagParsing = true

	cmd.AddCommand(NewUninstallAckCmd())
	cmd.AddCommand(NewUninstallAliasCmds(ackControllers, "ack-")...)
	cmd.AddCommand(adot.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallArgoCmd())
	cmd.AddCommand(NewUninstallAliasCmds(argoApps, "argo-")...)
	cmd.AddCommand(NewUninstallAutoscalingCmd())
	cmd.AddCommand(NewUninstallAliasCmds(autoscalingApps, "autoscaling-")...)
	cmd.AddCommand(NewUninstallAliasCmds(autoscalingApps, "as-")...)
	cmd.AddCommand(aws_fluent_bit.NewApp().NewUninstallCmd())
	cmd.AddCommand(awslbc.NewApp().NewUninstallCmd())
	cmd.AddCommand(certmanager.NewApp().NewUninstallCmd())
	cmd.AddCommand(cilium.NewApp().NewUninstallCmd())
	cmd.AddCommand(consul.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallContainerInsightsCmd())
	cmd.AddCommand(NewUninstallAliasCmds(containerInsightsApps, "container-insights-")...)
	cmd.AddCommand(NewUninstallAliasCmds(containerInsightsApps, "ci-")...)
	cmd.AddCommand(coredumphandler.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallCrossplaneCmd())
	cmd.AddCommand(NewUninstallAliasCmds(crossplane, "crossplane-")...)
	cmd.AddCommand(NewUninstallAliasCmds(crossplane, "cp-")...)
	cmd.AddCommand(NewUninstallExampleCmd())
	cmd.AddCommand(NewUninstallAliasCmds(exampleApps, "example-")...)
	cmd.AddCommand(NewUninstallAliasCmds(exampleApps, "ex-")...)
	cmd.AddCommand(externaldns.New().NewUninstallCmd())
	cmd.AddCommand(falco.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallFluxCmd())
	cmd.AddCommand(NewUninstallAliasCmds(fluxApps, "flux-")...)
	cmd.AddCommand(goldilocks.NewApp().NewUninstallCmd())
	cmd.AddCommand(harbor.NewApp().NewUninstallCmd())
	cmd.AddCommand(headlamp.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallIngressCmd())
	cmd.AddCommand(NewUninstallAliasCmds(ingressControllers, "ingress-")...)
	cmd.AddCommand(NewUninstallIstioCmd())
	cmd.AddCommand(NewUninstallAliasCmds(istioApps, "istio-")...)
	cmd.AddCommand(k8sgpt.NewApp().NewUninstallCmd())
	cmd.AddCommand(karpenter.NewApp().NewUninstallCmd())
	cmd.AddCommand(keycloak_amg.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallKubePrometheusCmd())
	cmd.AddCommand(NewUninstallAliasCmds(kubePrometheusApps, "kube-prometheus-")...)
	cmd.AddCommand(NewUninstallAliasCmds(kubePrometheusApps, "kprom-")...)
	cmd.AddCommand(kube_state_metrics.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallKubecostCmd())
	cmd.AddCommand(NewUninstallAliasCmds(kubecostApps, "kubecost-")...)
	cmd.AddCommand(NewUninstallLinkerdCmd())
	cmd.AddCommand(NewUninstallAliasCmds(linkerdApps, "linkerd-")...)
	cmd.AddCommand(metrics_server.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallNeuronCmd())
	cmd.AddCommand(NewUninstallAliasCmds(neuron, "neuron-")...)
	cmd.AddCommand(NewUninstallPolicyCmd())
	cmd.AddCommand(NewUninstallAliasCmds(policyApps, "policy-")...)
	cmd.AddCommand(prometheus_node_exporter.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallSecretsCmd())
	cmd.AddCommand(NewUninstallAliasCmds(secrets, "secret-")...)
	cmd.AddCommand(NewUninstallAliasCmds(secrets, "secrets-")...)
	cmd.AddCommand(sparkoperator.NewApp().NewUninstallCmd())
	cmd.AddCommand(NewUninstallStorageCmd())
	cmd.AddCommand(NewUninstallAliasCmds(storageApps, "storage-")...)
	cmd.AddCommand(vault.NewApp().NewUninstallCmd())
	cmd.AddCommand(velero.NewApp().NewUninstallCmd())
	cmd.AddCommand(vpclattice.NewApp().NewUninstallCmd())

	// Hidden commands for popular apps without using the group
	cmd.AddCommand(NewUninstallAliasCmds([]func() *application.Application{cd.NewApp}, "argo")...)
	cmd.AddCommand(NewUninstallAliasCmds([]func() *application.Application{cluster_autoscaler.NewApp}, "")...)
	cmd.AddCommand(NewUninstallAliasCmds([]func() *application.Application{ebs_csi.NewApp}, "")...)

	return cmd
}

// This creates alias commands for subcommands under INSTALL
func NewUninstallAliasCmds(appList []func() *application.Application, prefix string) []*cobra.Command {
	cmds := make([]*cobra.Command, 0, len(appList))

	for _, app := range appList {
		a := app()
		a.Command.Name = prefix + a.Command.Name
		a.Command.Hidden = true
		for i, alias := range a.Command.Aliases {
			a.Command.Aliases[i] = prefix + alias
		}
		cmds = append(cmds, a.NewUninstallCmd())
	}

	return cmds
}
