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
// Version: Latest is Chart/App 1.103.3 (as of 5/24/23)

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
				LatestChart:   "1.103.3",
				Latest:        "1.103.3",
				PreviousChart: "1.100.0",
				Previous:      "1.100.0",
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
serviceAccount:
  name: {{ .ServiceAccount }}
kubecostProductConfigs:
  clusterName: {{ .ClusterName }}
`
