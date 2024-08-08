package ascp

// https://github.com/aws/secrets-store-csi-driver-provider-aws/blob/main/examples/ExampleDeployment.yaml
const secretsProviderClassTemplate = `---
apiVersion: secrets-store.csi.x-k8s.io/v1
kind: SecretProviderClass
metadata:
  name: nginx-deployment-aws-secrets
spec:
  provider: aws
  parameters:
    objects: |
        - objectName: "MySecret"
          objectType: "secretsmanager"
`

const serviceAccountTemplate = `---
apiVersion: v1
kind: ServiceAccount
metadata:
  annotations:
    {{ .IrsaAnnotation }}
  name: {{ .ServiceAccount}}
  namespace: {{ .Namespace }}
`

// https://github.com/aws/secrets-store-csi-driver-provider-aws/blob/main/examples/ExampleDeployment.yaml
const serviceAndDeploymentTemplate = `---
kind: Service
apiVersion: v1
metadata:
  name: nginx-deployment
  namespace: {{ .Namespace }}
  labels:
    app: nginx
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
      targetPort: 80
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: {{ .Namespace }}
  labels:
    app: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      serviceAccountName: nginx-deployment-sa
      volumes:
      - name: secrets-store-inline
        csi:
          driver: secrets-store.csi.k8s.io
          readOnly: true
          volumeAttributes:
            secretProviderClass: "nginx-deployment-aws-secrets"
      containers:
      - name: nginx-deployment
        image: nginx
        ports:
        - containerPort: 80
        volumeMounts:
        - name: secrets-store-inline
          mountPath: "/mnt/secrets-store"
          readOnly: true
`
