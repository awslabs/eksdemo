package deviceplugin

// Manifest: https://github.com/aws-neuron/aws-neuron-sdk/blob/master/src/k8/k8s-neuron-device-plugin.yml
const daemonsetTemplate = `---
# https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: neuron-device-plugin-daemonset
  namespace: {{ .Namespace }}
spec:
  selector:
    matchLabels:
      name:  neuron-device-plugin-ds
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        name: neuron-device-plugin-ds
    spec:
      serviceAccount: neuron-device-plugin
      tolerations:
      - key: CriticalAddonsOnly
        operator: Exists
      - key: aws.amazon.com/neuron
        operator: Exists
        effect: NoSchedule
      # Mark this pod as a critical add-on; when enabled, the critical add-on
      # scheduler reserves resources for critical add-on pods so that they can
      # be rescheduled after a failure.
      # See https://kubernetes.io/docs/tasks/administer-cluster/guaranteed-scheduling-critical-addon-pods/
      priorityClassName: "system-node-critical"
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
              - matchExpressions:
                  - key: "node.kubernetes.io/instance-type"
                    operator: In
                    values:
                      - inf1.xlarge
                      - inf1.2xlarge
                      - inf1.6xlarge
                      - inf1.24xlarge
                      - inf2.xlarge
                      - inf2.4xlarge
                      - inf2.8xlarge
                      - inf2.24xlarge
                      - inf2.48xlarge
                      - trn1.2xlarge
                      - trn1.32xlarge
                      - trn1n.32xlarge
      containers:
        #Device Plugin containers are available both in us-east and us-west ecr
        #repos
      - image: public.ecr.aws/neuron/neuron-device-plugin:{{ or .Version "latest" }}
        imagePullPolicy: Always
        name: neuron-device-plugin
        env:
        - name: KUBECONFIG
          value: /etc/kubernetes/kubelet.conf
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop: ["ALL"]
        volumeMounts:
          - name: device-plugin
            mountPath: /var/lib/kubelet/device-plugins
          - name: infa-map
            mountPath: /run
      volumes:
        - name: device-plugin
          hostPath:
            path: /var/lib/kubelet/device-plugins
        - name: infa-map
          hostPath:
            path: /run
`

// Manifest: https://raw.githubusercontent.com/aws-neuron/aws-neuron-sdk/master/src/k8/k8s-neuron-device-plugin-rbac.yml
const rbacTemplate = `---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: neuron-device-plugin
rules:
- apiGroups:
  - ""
  resources:
  - nodes
  verbs:
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
  - patch
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - update
  - patch
  - get
  - list
  - watch
- apiGroups:
  - ""
  resources:
  - nodes/status
  verbs:
  - patch
  - update
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ .ServiceAccount }}
  namespace: {{ .Namespace }}
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: neuron-device-plugin
  namespace: {{ .Namespace }}
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: neuron-device-plugin
subjects:
- kind: ServiceAccount
  name: {{ .ServiceAccount }}
  namespace: {{ .Namespace }}
`
