package kube_ops_view

// https://codeberg.org/hjacobs/kube-ops-view/src/branch/main/deploy/deployment.yaml
const deploymentTemplate = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    application: kube-ops-view
    component: frontend
  name: kube-ops-view
  namespace: {{ .Namespace }}
spec:
  replicas: 1
  selector:
    matchLabels:
      application: kube-ops-view
      component: frontend
  template:
    metadata:
      labels:
        application: kube-ops-view
        component: frontend
    spec:
      serviceAccountName: kube-ops-view
      containers:
      - name: service
        # see https://github.com/hjacobs/kube-ops-view/releases
        image: hjacobs/kube-ops-view:{{ .Version }}
        args:
        # remove this option to use built-in memory store
        # - --redis-url=redis://kube-ops-view-redis:6379
        # example to add external links for nodes and pods
        # - --node-link-url-template=https://kube-web-view.example.org/clusters/{cluster}/nodes/{name}
        # - --pod-link-url-template=https://kube-web-view.example.org/clusters/{cluster}/namespaces/{namespace}/pods/{name}
        ports:
        - containerPort: 8080
          protocol: TCP
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          timeoutSeconds: 1
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 30
          timeoutSeconds: 10
          failureThreshold: 5
        resources:
          limits:
            cpu: 200m
            memory: 100Mi
          requests:
            cpu: 50m
            memory: 50Mi
        securityContext:
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 1000
`

// https://codeberg.org/hjacobs/kube-ops-view/src/branch/main/deploy/rbac.yaml
const rbacTemplate = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-ops-view
  namespace: {{ .Namespace }}
---
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-ops-view
rules:
- apiGroups: [""]
  resources: ["nodes", "pods"]
  verbs:
    - list
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes", "pods"]
  verbs:
    - get
    - list
---
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kube-ops-view
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-ops-view
subjects:
- kind: ServiceAccount
  name: kube-ops-view
  namespace: {{ .Namespace }}
`

// https://codeberg.org/hjacobs/kube-ops-view/src/branch/main/deploy/service.yaml
const serviceTemplate = `---
apiVersion: v1
kind: Service
metadata:
  labels:
    application: kube-ops-view
    component: frontend
  name: kube-ops-view
  namespace: {{ .Namespace }}
  annotations:
    {{- .ServiceAnnotations | nindent 4 }}
spec:
  selector:
    application: kube-ops-view
    component: frontend
  type: {{ .ServiceType }}
  ports:
  - port: 80
    protocol: TCP
    targetPort: 8080
{{- if .IngressHost }}
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
  name: kube-ops-view
  namespace: {{ .Namespace }}
spec:
  ingressClassName: {{ .IngressClass }}
  rules:
    - host: {{ .IngressHost }}
      http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: kube-ops-view
              port:
                number: 80
  tls:
  - hosts:
    - {{ .IngressHost }}
  {{- if ne .IngressClass "alb" }}
    secretName: kube-ops-view-tls
  {{- end }}
{{- end }}
`
