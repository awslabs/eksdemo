package cloudwatch_agent

const kustomizeTemplate = `---
resources:
- manifest.yaml
patches:
- patch: |-
    - op: add
      path: /metadata/annotations
      value:
        {{ .IrsaAnnotation }}
  target:
    kind: ServiceAccount
    name: {{ .ServiceAccount }}
    namespace: {{ .Namespace }}
    version: v1
`

// Manfiest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cwagent/cwagent-serviceaccount.yaml
const serviceAccountTemplate = `---
# create cwagent service account and role binding
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .ServiceAccount }}
  namespace: {{ .Namespace }}

---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cloudwatch-agent-role
rules:
  - apiGroups: [""]
    resources: ["pods", "nodes", "endpoints"]
    verbs: ["list", "watch"]
  - apiGroups: ["apps"]
    resources: ["replicasets"]
    verbs: ["list", "watch"]
  - apiGroups: ["batch"]
    resources: ["jobs"]
    verbs: ["list", "watch"]
  - apiGroups: [""]
    resources: ["nodes/proxy"]
    verbs: ["get"]
  - apiGroups: [""]
    resources: ["nodes/stats", "configmaps", "events"]
    verbs: ["create"]
  - apiGroups: [""]
    resources: ["configmaps"]
    resourceNames: ["cwagent-clusterleader"]
    verbs: ["get","update"]

---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cloudwatch-agent-role-binding
subjects:
  - kind: ServiceAccount
    name: {{ .ServiceAccount }}
    namespace: {{ .Namespace }}
roleRef:
  kind: ClusterRole
  name: cloudwatch-agent-role
  apiGroup: rbac.authorization.k8s.io
`

// Manfiest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cwagent/cwagent-configmap.yaml
// Note: Only difference from the manifest above is the addition of the agent.region. That is included in the Quick Start config below:
// https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/quickstart/cwagent-fluent-bit-quickstart.yaml#L65
const configMapTemplate = `---
# create configmap for cwagent config
apiVersion: v1
data:
  # Configuration is in Json format. No matter what configure change you make,
  # please keep the Json blob valid.
  cwagentconfig.json: |
    {
      "agent": {
        "region": "{{ .Region }}"
      },
      "logs": {
        "metrics_collected": {
          "kubernetes": {
            "cluster_name": "{{ .ClusterName }}",
            "metrics_collection_interval": 60
          }
        },
        "force_flush_interval": 5
      }
    }
kind: ConfigMap
metadata:
  name: cwagentconfig
  namespace: {{ .Namespace }}
`

// Manifest: https://github.com/aws-samples/amazon-cloudwatch-container-insights/blob/master/k8s-deployment-manifest-templates/deployment-mode/daemonset/container-insights-monitoring/cwagent/cwagent-daemonset.yaml
const daemonSetTemplate = `---
# deploy cwagent as daemonset
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: cloudwatch-agent
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      name: cloudwatch-agent
  template:
    metadata:
      labels:
        name: cloudwatch-agent
    spec:
      containers:
        - name: cloudwatch-agent
          image: amazon/cloudwatch-agent:{{ or .Version "latest" }}
          #ports:
          #  - containerPort: 8125
          #    hostPort: 8125
          #    protocol: UDP
          resources:
            limits:
              cpu:  200m
              memory: 200Mi
            requests:
              cpu: 200m
              memory: 200Mi
          # Please don't change below envs
          env:
            - name: HOST_IP
              valueFrom:
                fieldRef:
                  fieldPath: status.hostIP
            - name: HOST_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
            - name: K8S_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: CI_VERSION
              value: "k8s/1.3.10"
          # Please don't change the mountPath
          volumeMounts:
            - name: cwagentconfig
              mountPath: /etc/cwagentconfig
            - name: rootfs
              mountPath: /rootfs
              readOnly: true
            - name: dockersock
              mountPath: /var/run/docker.sock
              readOnly: true
            - name: varlibdocker
              mountPath: /var/lib/docker
              readOnly: true
            - name: containerdsock
              mountPath: /run/containerd/containerd.sock
              readOnly: true
            - name: sys
              mountPath: /sys
              readOnly: true
            - name: devdisk
              mountPath: /dev/disk
              readOnly: true
      volumes:
        - name: cwagentconfig
          configMap:
            name: cwagentconfig
        - name: rootfs
          hostPath:
            path: /
        - name: dockersock
          hostPath:
            path: /var/run/docker.sock
        - name: varlibdocker
          hostPath:
            path: /var/lib/docker
        - name: containerdsock
          hostPath:
            path: /run/containerd/containerd.sock
        - name: sys
          hostPath:
            path: /sys
        - name: devdisk
          hostPath:
            path: /dev/disk/
      terminationGracePeriodSeconds: 60
      serviceAccountName: cloudwatch-agent
`
