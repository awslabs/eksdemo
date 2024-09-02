package karpenter

// JSON: https://github.com/aws/karpenter-provider-aws/blob/main/website/content/en/preview/getting-started/getting-started-with-karpenter/karpenter-performance-dashboard.json
//
//	kubectl create configmap karpenter-performance \
//	  --from-file=./karpenter-performance-dashboard.json \
//	  --dry-run=client -o yaml > performance_dashboard.yaml
//
// Must double escape because textTemplate transforms and Helm uses the same Golang template
// TextTemplate turns: "{{ "{{" }} {{ "{{controller}}" | printf "%q" }} {{ "}}" }}" _to_ "{{ "{{controller}}" }}"
// and Helm turns it _to_ {{controller}}
const performanceDashboardTemplate = `---
apiVersion: v1
kind: ConfigMap
metadata:
  name: karpenter-performance
data:
  karpenter-performance-dashboard.json: |
    {
        "annotations": {
          "list": [
            {
              "builtIn": 1,
              "datasource": "-- Grafana --",
              "enable": true,
              "hide": true,
              "iconColor": "rgba(0, 211, 255, 1)",
              "name": "Annotations & Alerts",
              "target": {
                "limit": 100,
                "matchAny": false,
                "tags": [],
                "type": "dashboard"
              },
              "type": "dashboard"
            }
          ]
        },
        "editable": true,
        "fiscalYearStartMonth": 0,
        "graphTooltip": 0,
        "id": 33,
        "links": [],
        "liveNow": true,
        "panels": [
          {
            "datasource": {
              "type": "prometheus",
              "uid": "${datasource}"
            },
            "fieldConfig": {
              "defaults": {
                "color": {
                  "mode": "palette-classic"
                },
                "custom": {
                  "axisBorderShow": false,
                  "axisCenteredZero": false,
                  "axisColorMode": "text",
                  "axisLabel": "",
                  "axisPlacement": "auto",
                  "barAlignment": 0,
                  "drawStyle": "line",
                  "fillOpacity": 0,
                  "gradientMode": "none",
                  "hideFrom": {
                    "legend": false,
                    "tooltip": false,
                    "viz": false
                  },
                  "insertNulls": false,
                  "lineInterpolation": "linear",
                  "lineWidth": 1,
                  "pointSize": 5,
                  "scaleDistribution": {
                    "type": "linear"
                  },
                  "showPoints": "never",
                  "spanNulls": false,
                  "stacking": {
                    "group": "A",
                    "mode": "none"
                  },
                  "thresholdsStyle": {
                    "mode": "off"
                  }
                },
                "mappings": [],
                "thresholds": {
                  "mode": "absolute",
                  "steps": [
                    {
                      "color": "green",
                      "value": null
                    },
                    {
                      "color": "red",
                      "value": 80
                    }
                  ]
                },
                "unit": "s"
              },
              "overrides": []
            },
            "gridPos": {
              "h": 9,
              "w": 24,
              "x": 0,
              "y": 0
            },
            "id": 4,
            "options": {
              "legend": {
                "calcs": [
                  "max"
                ],
                "displayMode": "table",
                "placement": "right",
                "showLegend": true
              },
              "tooltip": {
                "mode": "single",
                "sort": "none"
              }
            },
            "targets": [
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_nodes_termination_duration_seconds{quantile=\"0\"})",
                "legendFormat": "Min",
                "range": true,
                "refId": "A"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_nodes_termination_duration_seconds{quantile=\"0.5\"})",
                "hide": false,
                "legendFormat": "P50",
                "range": true,
                "refId": "B"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_nodes_termination_duration_seconds{quantile=\"0.9\"})",
                "hide": false,
                "legendFormat": "P90",
                "range": true,
                "refId": "C"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_nodes_termination_duration_seconds{quantile=\"0.99\"})",
                "hide": false,
                "legendFormat": "P99",
                "range": true,
                "refId": "D"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_nodes_termination_duration_seconds{quantile=\"1\"})",
                "hide": false,
                "legendFormat": "Max",
                "range": true,
                "refId": "E"
              }
            ],
            "title": "Node Termination Latency",
            "type": "timeseries"
          },
          {
            "datasource": {
              "type": "prometheus",
              "uid": "${datasource}"
            },
            "fieldConfig": {
              "defaults": {
                "color": {
                  "mode": "palette-classic"
                },
                "custom": {
                  "axisBorderShow": false,
                  "axisCenteredZero": false,
                  "axisColorMode": "text",
                  "axisLabel": "",
                  "axisPlacement": "auto",
                  "barAlignment": 0,
                  "drawStyle": "line",
                  "fillOpacity": 10,
                  "gradientMode": "none",
                  "hideFrom": {
                    "legend": false,
                    "tooltip": false,
                    "viz": false
                  },
                  "insertNulls": false,
                  "lineInterpolation": "linear",
                  "lineWidth": 1,
                  "pointSize": 5,
                  "scaleDistribution": {
                    "type": "linear"
                  },
                  "showPoints": "never",
                  "spanNulls": false,
                  "stacking": {
                    "group": "A",
                    "mode": "none"
                  },
                  "thresholdsStyle": {
                    "mode": "off"
                  }
                },
                "mappings": [],
                "min": 0,
                "thresholds": {
                  "mode": "absolute",
                  "steps": [
                    {
                      "color": "green",
                      "value": null
                    },
                    {
                      "color": "red",
                      "value": 80
                    }
                  ]
                },
                "unit": "s"
              },
              "overrides": []
            },
            "gridPos": {
              "h": 8,
              "w": 24,
              "x": 0,
              "y": 9
            },
            "id": 2,
            "options": {
              "legend": {
                "calcs": [
                  "max"
                ],
                "displayMode": "table",
                "placement": "right",
                "showLegend": true
              },
              "tooltip": {
                "mode": "single",
                "sort": "none"
              }
            },
            "targets": [
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_pods_startup_duration_seconds{quantile=\"0\"})",
                "format": "time_series",
                "legendFormat": "Min",
                "range": true,
                "refId": "Minimum"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_pods_startup_duration_seconds{quantile=\"0.5\"})",
                "hide": false,
                "legendFormat": "P50",
                "range": true,
                "refId": "Median"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_pods_startup_duration_seconds{quantile=\"0.9\"})",
                "hide": false,
                "legendFormat": "P90",
                "range": true,
                "refId": "P90"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_pods_startup_duration_seconds{quantile=\"0.99\"})",
                "hide": false,
                "legendFormat": "P99",
                "range": true,
                "refId": "P99"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "max(karpenter_pods_startup_duration_seconds{quantile=\"1\"})",
                "hide": false,
                "legendFormat": "Max",
                "range": true,
                "refId": "Maximum"
              }
            ],
            "title": "Pod Startup Latency",
            "type": "timeseries"
          },
          {
            "datasource": {
              "type": "prometheus",
              "uid": "${datasource}"
            },
            "fieldConfig": {
              "defaults": {
                "color": {
                  "mode": "palette-classic"
                },
                "custom": {
                  "axisBorderShow": false,
                  "axisCenteredZero": false,
                  "axisColorMode": "text",
                  "axisLabel": "",
                  "axisPlacement": "auto",
                  "barAlignment": 0,
                  "drawStyle": "line",
                  "fillOpacity": 10,
                  "gradientMode": "none",
                  "hideFrom": {
                    "legend": false,
                    "tooltip": false,
                    "viz": false
                  },
                  "insertNulls": false,
                  "lineInterpolation": "linear",
                  "lineWidth": 1,
                  "pointSize": 5,
                  "scaleDistribution": {
                    "type": "linear"
                  },
                  "showPoints": "never",
                  "spanNulls": false,
                  "stacking": {
                    "group": "A",
                    "mode": "none"
                  },
                  "thresholdsStyle": {
                    "mode": "off"
                  }
                },
                "mappings": [],
                "thresholds": {
                  "mode": "absolute",
                  "steps": [
                    {
                      "color": "green",
                      "value": null
                    },
                    {
                      "color": "red",
                      "value": 80
                    }
                  ]
                },
                "unit": "s"
              },
              "overrides": []
            },
            "gridPos": {
              "h": 8,
              "w": 13,
              "x": 0,
              "y": 17
            },
            "id": 6,
            "options": {
              "legend": {
                "calcs": [],
                "displayMode": "list",
                "placement": "bottom",
                "showLegend": true
              },
              "tooltip": {
                "mode": "single",
                "sort": "none"
              }
            },
            "targets": [
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "histogram_quantile(0, rate(controller_runtime_reconcile_time_seconds_bucket{controller=\"$controller\"}[10m]))",
                "hide": false,
                "legendFormat": "Min",
                "range": true,
                "refId": "Minimum"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "histogram_quantile(0.5, rate(controller_runtime_reconcile_time_seconds_bucket{controller=\"$controller\"}[10m]))",
                "legendFormat": "P50",
                "range": true,
                "refId": "Median"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "histogram_quantile(0.9, rate(controller_runtime_reconcile_time_seconds_bucket{controller=\"$controller\"}[10m]))",
                "hide": false,
                "legendFormat": "P90",
                "range": true,
                "refId": "P90"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "histogram_quantile(0.99, rate(controller_runtime_reconcile_time_seconds_bucket{controller=\"$controller\"}[10m]))",
                "hide": false,
                "legendFormat": "P99",
                "range": true,
                "refId": "P99"
              },
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "histogram_quantile(1, rate(controller_runtime_reconcile_time_seconds_bucket{controller=\"$controller\"}[10m]))",
                "hide": false,
                "legendFormat": "Max",
                "range": true,
                "refId": "Maximum"
              }
            ],
            "title": "Controller Reconciliation Latency [$controller]",
            "type": "timeseries"
          },
          {
            "datasource": {
              "type": "prometheus",
              "uid": "${datasource}"
            },
            "fieldConfig": {
              "defaults": {
                "color": {
                  "mode": "thresholds"
                },
                "mappings": [],
                "thresholds": {
                  "mode": "absolute",
                  "steps": [
                    {
                      "color": "green",
                      "value": null
                    },
                    {
                      "color": "red",
                      "value": 80
                    }
                  ]
                },
                "unit": "reqps"
              },
              "overrides": []
            },
            "gridPos": {
              "h": 8,
              "w": 11,
              "x": 13,
              "y": 17
            },
            "id": 8,
            "options": {
              "displayMode": "gradient",
              "maxVizHeight": 300,
              "minVizHeight": 10,
              "minVizWidth": 0,
              "namePlacement": "auto",
              "orientation": "horizontal",
              "reduceOptions": {
                "calcs": [
                  "lastNotNull"
                ],
                "fields": "",
                "values": false
              },
              "showUnfilled": true,
              "sizing": "auto",
              "valueMode": "color"
            },
            "pluginVersion": "11.1.0",
            "targets": [
              {
                "datasource": {
                  "type": "prometheus",
                  "uid": "${datasource}"
                },
                "editorMode": "code",
                "expr": "sum(rate(controller_runtime_reconcile_total{job=\"karpenter\"}[10m])) by (controller)",
                "legendFormat": "{{ "{{" }} {{ "{{controller}}" | printf "%q" }} {{ "}}" }}",
                "range": true,
                "refId": "A"
              }
            ],
            "title": "Controller Reconciliation Rate",
            "type": "bargauge"
          }
        ],
        "refresh": "5s",
        "schemaVersion": 39,
        "tags": [],
        "templating": {
          "list": [
            {
              "current": {
                "selected": false,
                "text": "Prometheus",
                "value": "prometheus"
              },
              "hide": 0,
              "includeAll": false,
              "label": "Data Source",
              "multi": false,
              "name": "datasource",
              "options": [],
              "query": "prometheus",
              "refresh": 1,
              "regex": "",
              "skipUrlSync": false,
              "type": "datasource"
            },
            {
              "current": {
                "selected": false,
                "text": "disruption",
                "value": "disruption"
              },
              "datasource": {
                "type": "prometheus",
                "uid": "${datasource}"
              },
              "definition": "label_values(controller_runtime_reconcile_time_seconds_count{job=\"karpenter\"}, controller)",
              "hide": 0,
              "includeAll": false,
              "multi": false,
              "name": "controller",
              "options": [],
              "query": {
                "query": "label_values(controller_runtime_reconcile_time_seconds_count{job=\"karpenter\"}, controller)",
                "refId": "StandardVariableQuery"
              },
              "refresh": 2,
              "regex": "",
              "skipUrlSync": false,
              "sort": 0,
              "type": "query"
            }
          ]
        },
        "time": {
          "from": "now-3h",
          "to": "now"
        },
        "timepicker": {},
        "timezone": "",
        "title": "Karpenter Performance v1",
        "uid": "fdusq1f2alerke",
        "version": 3,
        "weekStart": ""
      }
`
