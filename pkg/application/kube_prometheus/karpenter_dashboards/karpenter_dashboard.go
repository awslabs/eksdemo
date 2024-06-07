package karpenter_dashboards

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:       https://karpenter.sh/docs/
// GitHub:     https://github.com/aws/karpenter-provider-aws
// Dashboards: https://github.com/aws/karpenter-provider-aws/tree/main/website/content/en/preview/getting-started/getting-started-with-karpenter

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "kube-prometheus",
			Name:        "karpenter-dashboards",
			Description: "Karpenter Dashboards and ServiceMonitor",
			Aliases:     []string{"karpenter"},
		},

		Flags: flags,

		Installer: &installer.ManifestInstaller{
			AppName: "kube-prometheus-karpenter-dashboards",
			ResourceTemplate: &template.TextTemplate{
				Template: serviceMonitorTemplate + capacityDashboardTemplate + performanceDashboardTemplate,
			},
			KustomizeTemplate: &template.TextTemplate{
				Template: kustomizeTemplate,
			},
		},

		Options: options,
	}
}

const kustomizeTemplate = `---
resources:
- manifest.yaml
patches:
# Adds label to ConfigMaps to match grafana-sc-dashboard sidecar
- patch: |-
    - op: add
      path: /metadata/labels
      value:
        grafana_dashboard: "1"
  target:
    kind: ConfigMap
    version: v1
`
