package otel

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

// As of 9/8/24, this doesn't work due to https://github.com/aws-observability/aws-otel-collector/issues/2470
func NewSimplestCollector() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "simplest",
			Description: "Simplest Collector",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: simplestCollectorTemplate,
			},
		},

		Options: &resource.CommonOptions{
			Name:          "simplest",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}
}

// https://github.com/open-telemetry/opentelemetry-operator#getting-started
const simplestCollectorTemplate = `---
apiVersion: opentelemetry.io/v1beta1
kind: OpenTelemetryCollector
metadata:
  namespace: {{ .Namespace }}
  name: {{ .Name }}
spec:
  config:
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
    processors:
      memory_limiter:
        check_interval: 1s
        limit_percentage: 75
        spike_limit_percentage: 15
      batch:
        send_batch_size: 10000
        timeout: 10s

    exporters:
      debug: {}

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: [memory_limiter, batch]
          exporters: [debug]
`
