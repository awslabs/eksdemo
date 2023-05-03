package istiod

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

// https://github.com/istio/istio/blob/master/samples/addons/extras/prometheus-operator.yaml

func serviceMonitor() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "istiod-service-monitor",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: serviceMonitorYamlTemplate,
			},
		},
	}
	return res
}

const serviceMonitorYamlTemplate = `---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: istio-component-monitor
  namespace: istio-system
  labels:
    monitoring: istio-components
    release: istio
spec:
  jobLabel: istio
  targetLabels: [app]
  selector:
    matchExpressions:
    - {key: istio, operator: In, values: [pilot]}
  namespaceSelector:
    any: true
  endpoints:
  - port: http-monitoring
    interval: 15s
`
