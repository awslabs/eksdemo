package kubecost_eks

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://docs.aws.amazon.com/eks/latest/userguide/cost-monitoring.html
// Docs:    https://guide.kubecost.com/
// Helm:    https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer
// Values:  https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values.yaml
// Values:  https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-eks-cost-monitoring.yaml
// Repo:    https://gallery.ecr.aws/kubecost/cost-model
// Repo:    https://gallery.ecr.aws/kubecost/frontend
// Repo:    https://gallery.ecr.aws/kubecost/prometheus
// Repo:    https://gallery.ecr.aws/bitnami/configmap-reload
// Version: Latest is Chart/App 1.103.3 (as of 5/24/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kubecost",
			Name:        "eks",
			Description: "EKS optimized bundle of Kubecost",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost-eks",
			RepositoryURL: "oci://public.ecr.aws/kubecost/cost-analyzer",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app
}

const valuesTemplate = `---
fullnameOverride: kubecost-cost-analyzer
global:
  grafana:
    # If false, Grafana will not be installed
    enabled: false
    # If true, the kubecost frontend will route to your grafana through its service endpoint
    proxy: false
podSecurityPolicy:
  enabled: false
imageVersion: prod-{{ .Version }}
kubecostFrontend:
  image: public.ecr.aws/kubecost/frontend
kubecostModel:
  image: public.ecr.aws/kubecost/cost-model
  # The total number of days the ETL storage will build
  etlStoreDurationDays: 120
serviceAccount:
  name: {{ .ServiceAccount }}
{{- if .IngressHost }}
ingress:
  enabled: true
  className: {{ .IngressClass }}
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  pathType: Prefix
  hosts:
    - {{ .IngressHost }}
  tls:
  - hosts:
    - {{ .IngressHost }}
  {{- if ne .IngressClass "alb" }}
    secretName: cost-analyzer-tls
  {{- end }}
{{- end }}
service:
  type: {{ .ServiceType }}
  annotations: 
    {{- .ServiceAnnotations | nindent 4 }}
prometheus:
  kube-state-metrics:
    fullnameOverride: kubecost-kube-state-metrics
  server:
    fullnameOverride: kubecost-prometheus-server
    image:
      repository: public.ecr.aws/kubecost/prometheus
      tag: v2.35.0
    global:
      # overrides kubecost default of 60s, sets to prom default of 10s
      scrape_timeout: 10s
      external_labels:
        cluster_id: {{ .ClusterName }} # Each cluster should have a unique ID
  configmapReload:
    prometheus:
      enabled: false
      image:
        repository: public.ecr.aws/bitnami/configmap-reload
        tag: 0.7.1
  nodeExporter:
    enabled: {{ .EnableNodeExporter }}
    fullnameOverride: kubecost-prometheus-node-exporter
reporting:
  productAnalytics: false
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`
