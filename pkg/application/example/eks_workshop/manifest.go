package eks_workshop

const manifestTemplate = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ecsdemo-nodejs
  labels:
    app: ecsdemo-nodejs
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .NodeReplicas }}
  selector:
    matchLabels:
      app: ecsdemo-nodejs
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: ecsdemo-nodejs
    spec:
      containers:
      - image: public.ecr.aws/aws-containers/ecsdemo-nodejs:latest
        imagePullPolicy: Always
        name: ecsdemo-nodejs
        ports:
        - containerPort: 3000
          protocol: TCP
---
apiVersion: v1
kind: Service
metadata:
  name: ecsdemo-nodejs
  namespace: {{ .Namespace }}
spec:
  selector:
    app: ecsdemo-nodejs
  ports:
   -  protocol: TCP
      port: 80
      targetPort: 3000
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ecsdemo-crystal
  labels:
    app: ecsdemo-crystal
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .CrystalReplicas }}
  selector:
    matchLabels:
      app: ecsdemo-crystal
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: ecsdemo-crystal
    spec:
      containers:
      - image: public.ecr.aws/aws-containers/ecsdemo-crystal:latest
        imagePullPolicy: Always
        name: ecsdemo-crystal
        ports:
        - containerPort: 3000
          protocol: TCP
---
apiVersion: v1
kind: Service
metadata:
  name: ecsdemo-crystal
  namespace: {{ .Namespace }}
spec:
  selector:
    app: ecsdemo-crystal
  ports:
   -  protocol: TCP
      port: 80
      targetPort: 3000
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ecsdemo-frontend
  labels:
    app: ecsdemo-frontend
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .FrontendReplicas }}
  selector:
    matchLabels:
      app: ecsdemo-frontend
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: ecsdemo-frontend
    spec:
      containers:
      - image: public.ecr.aws/aws-containers/ecsdemo-frontend:latest
        imagePullPolicy: Always
        name: ecsdemo-frontend
        ports:
        - containerPort: 3000
          protocol: TCP
        env:
        - name: CRYSTAL_URL
          value: "http://ecsdemo-crystal.{{ .Namespace }}.svc.cluster.local/crystal"
        - name: NODEJS_URL
          value: "http://ecsdemo-nodejs.{{ .Namespace }}.svc.cluster.local/"

---
apiVersion: v1
kind: Service
metadata:
  name: ecsdemo-frontend
  annotations:
    {{- .ServiceAnnotations | nindent 4 }}
spec:
  selector:
    app: ecsdemo-frontend
  type: {{ .ServiceType }}
  ports:
   -  protocol: TCP
      port: 80
      targetPort: 3000
{{- if .IngressHost }}
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  namespace: {{ .Namespace }}
  name: ecsdemo-frontend
  annotations:
    {{- .IngressAnnotations | nindent 4 }}
spec:
  ingressClassName: {{ .IngressClass }}
  rules:
  - http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: ecsdemo-frontend
            port:
              number: 80
  tls:
  - hosts:
    - {{ .IngressHost }}
  {{- if ne .IngressClass "alb" }}
    secretName: eks-workshop-tls
  {{- end }}
{{- end }}
`
