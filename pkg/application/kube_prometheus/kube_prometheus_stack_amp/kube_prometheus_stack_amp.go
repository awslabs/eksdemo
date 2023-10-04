package kube_prometheus_stack_amp

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://github.com/prometheus-operator/kube-prometheus/tree/main/docs
// GitHub:  https://github.com/prometheus-operator/kube-prometheus
// Helm:    https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack
// Helm:    https://github.com/grafana/helm-charts/tree/main/charts/grafana
// Repo:    https://quay.io/prometheus-operator/prometheus-operator
// Version: Latest is Chart 51.2.0, PromOperator v0.68.0 (as of 10/3/23)

func NewApp() *application.Application {
	options, flags := NewOptions()
	options.AmpWorkspaceOptions = &amp_workspace.AmpWorkspaceOptions{
		CommonOptions: resource.CommonOptions{
			Name: "amazon-managed-prometheus-workspace",
		},
	}

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kube-prometheus",
			Name:        "stack-amp",
			Description: "Kube Prometheus Stack using Amazon Managed Prometheus",
			Aliases:     []string{"amp"},
		},

		Dependencies: []*resource.Resource{
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "prometheus-amp-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: irsaPolicyDocument,
				},
			}),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name:           "grafana-amp-irsa",
					ServiceAccount: options.GrafanaServiceAccount,
				},
				PolicyType: irsa.PolicyARNs,
				Policy: []string{
					"arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess",
				},
			}),
			amp_workspace.NewResourceWithOptions(options.AmpWorkspaceOptions),
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "kube-prometheus-stack",
			ReleaseName:   "kube-prometheus-stack-amp",
			RepositoryURL: "https://prometheus-community.github.io/helm-charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			PostRenderKustomize: &template.TextTemplate{
				Template: postRenderKustomize,
			},
		},
	}

	app.Options = options
	app.Flags = flags

	return app
}

// https://docs.aws.amazon.com/prometheus/latest/userguide/set-up-irsa.html#set-up-irsa-ingest
const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Effect: Allow
  Action:
  - aps:RemoteWrite
  - aps:GetSeries
  - aps:GetLabels
  - aps:GetMetricMetadata
  Resource: "*"
`

const valuesTemplate = `---
fullnameOverride: prometheus
defaultRules:
  rules:
    alertmanager: false
alertmanager:
  enabled: false
grafana:
{{- if .DisableGrafana }}
  enabled: false
{{- else }}
  fullnameOverride: grafana
  serviceAccount:
    name: {{ .GrafanaServiceAccount }}
    annotations:
      {{ .IrsaAnnotationFor .GrafanaServiceAccount }}
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
  sidecar:
    datasources:
      defaultDatasourceEnabled: false
  additionalDataSources:
  - name: Amazon Managed Service for Prometheus
    type: prometheus
    url: {{ .AmpEndpoint }}
    access: proxy
    isDefault: true
    jsonData:
      sigV4Auth: true
      sigV4AuthType: default
      sigV4Region: {{ .Region }}
    timeInterval: 30s
  service:
    annotations:
      {{- .ServiceAnnotations | nindent 6 }}
    type: {{ .ServiceType }}
  # Temporary fix for issue: https://github.com/prometheus-community/helm-charts/issues/1867
  # Note the "serviceMonitorSelectorNilUsesHelmValues: false" below also resolves the issue
  serviceMonitor:
    labels:
      release: kube-prometheus-stack-amp
  adminPassword: {{ .GrafanaAdminPassword }}
  grafana.ini:
    auth:
      sigv4_auth_enabled: true
{{- end }}
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
  serviceAccount:
    name: {{ .ServiceAccount }}
    annotations:
      {{ .IrsaAnnotation }}
  prometheusSpec:
    # selects ServiceMonitors without the "release: kube-prometheus-stack-amp" label
    serviceMonitorSelectorNilUsesHelmValues: false
    remoteWrite:
    - url: {{ .AmpEndpoint }}api/v1/remote_write
      sigv4:
        region: {{ .Region }}
      queueConfig:
        maxSamplesPerSend: 1000
        maxShards: 200
        capacity: 2500
    remoteWriteDashboards: true
cleanPrometheusOperatorObjectNames: true
`

// https://github.com/kubernetes-sigs/kustomize/blob/master/examples/patchMultipleObjects.md
const postRenderKustomize = `---
resources:
- manifest.yaml
patches:
# Delete the AlertManager dashboard as it's disabled
- patch: |-
    $patch: delete
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: needed-but-not-used
  target:
    name: ".*alertmanager-overview$"
`
