package kube_prometheus_stack

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/prometheus-operator/kube-prometheus/tree/main/docs
// GitHub:  https://github.com/prometheus-operator/kube-prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Helm:    https://github.com/grafana/helm-charts/tree/main/charts/grafana
// Repo:    https://quay.io/prometheus-operator/prometheus-operator
// Version: Latest is Chart 51.2.0, PromOperator v0.68.0 (as of 10/3/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kube-prometheus",
			Name:        "stack",
			Description: "Kube Prometheus Stack",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kube-prometheus-stack",
			ReleaseName:   "kube-prometheus-stack",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
fullnameOverride: prometheus
grafana:
  fullnameOverride: grafana
  # Temporary fix for issue: https://github.com/prometheus-community/helm-charts/issues/1867
  # Note the "serviceMonitorSelectorNilUsesHelmValues: false" below also resolves the issue
  serviceMonitor:
    labels:
      release: kube-prometheus-stack
  adminPassword: {{ .GrafanaAdminPassword }}
{{- if .IngressHost }}
  ingress:
    enabled: true
    ingressClassName: {{ .IngressClass }}
    annotations:
      {{- .IngressAnnotations | nindent 6 }}
    hosts:
    - {{ .IngressHost }}
    tls:
    - hosts:
      - {{ .IngressHost }}
    {{- if ne .IngressClass "alb" }}
      secretName: grafana-tls
    {{- end}}
{{- end }}
  service:
    annotations:
      {{- .ServiceAnnotations | nindent 6 }}
    type: {{ .ServiceType }}
kubeControllerManager:
  enabled: false
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kube-state-metrics:
  fullnameOverride: kube-state-metrics
prometheus-node-exporter:
  fullnameOverride: node-exporter
prometheusOperator:
  image:
    tag: {{ .Version }}
prometheus:
  prometheusSpec:
    # selects ServiceMonitors without the "release: kube-prometheus-stack" label
    serviceMonitorSelectorNilUsesHelmValues: false
cleanPrometheusOperatorObjectNames: true
`
