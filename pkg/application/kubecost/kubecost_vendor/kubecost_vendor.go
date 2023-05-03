package kubecost_vendor

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://guide.kubecost.com/
// Helm:    https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer
// Repo:    gcr.io/kubecost1/cost-model
// Version: Latest is Chart/App 1.100.0 (as of 2/9/23)
// TODO:    remove podSecurityPolicy from values.yaml next update

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "kubecost",
			Name:        "vendor",
			Description: "Vendor distribution of Kubecost",
		},

		Options: &application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kubecost",
			ServiceAccount:               "kubecost-cost-analyzer",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.100.0",
				Latest:        "1.100.0",
				PreviousChart: "1.98.0",
				Previous:      "1.98.0",
			},
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "cost-analyzer",
			ReleaseName:   "kubecost-vendor",
			RepositoryURL: "https://kubecost.github.io/cost-analyzer/",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	return app
}

const valuesTemplate = `---
fullnameOverride: kubecost-cost-analyzer
global:
  grafana:
    # Required due to fullnameOverride on grafana
    fqdn: kubecost-grafana.{{ .Namespace }}
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
podSecurityPolicy:
  enabled: false
service:
  type: {{ .ServiceType }}
  annotations: 
    {{- .ServiceAnnotations | nindent 4 }}
prometheus:
  kube-state-metrics:
    fullnameOverride: kubecost-kube-state-metrics
  nodeExporter:
    fullnameOverride: kubecost-prometheus-node-exporter
  server:
    fullnameOverride: kubecost-prometheus-server
    global:
      external_labels:
        cluster_id: {{ .ClusterName }} # Each cluster should have a unique ID
grafana:
  fullnameOverride: kubecost-grafana
  rbac:
    pspEnabled: false
serviceAccount:
  name: {{ .ServiceAccount }}
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`
