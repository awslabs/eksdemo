package karpenter_dashboards

// Manifest: https://github.com/aws/karpenter/blob/main/charts/karpenter/templates/servicemonitor.yaml
const serviceMonitorTemplate = `---
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: karpenter
spec:
  endpoints:
  - path: /metrics
    port: http-metrics
  namespaceSelector:
    matchNames:
    - {{ .KarpenterNamespace }}
  selector:
    matchLabels:
      app.kubernetes.io/instance: autoscaling-karpenter
      app.kubernetes.io/name: karpenter
`
