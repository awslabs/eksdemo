package sparkpi

// https://github.com/kubeflow/spark-operator/blob/master/examples/spark-pi.yaml
const manifestTemplate = `---
apiVersion: sparkoperator.k8s.io/v1beta2
kind: SparkApplication
metadata:
  name: spark-pi
  namespace: {{ .Namespace }}
spec:
  type: Scala
  mode: cluster
  image: spark:3.5.0
  imagePullPolicy: IfNotPresent
  mainClass: org.apache.spark.examples.SparkPi
  mainApplicationFile: local:///opt/spark/examples/jars/spark-examples_2.12-3.5.0.jar
  sparkVersion: 3.5.0
  driver:
    labels:
      version: 3.5.0
    cores: 1
    coreLimit: 1200m
    memory: 512m
    serviceAccount: {{ .ServiceAccount }}
  executor:
    labels:
      version: 3.5.0
    instances: 1
    cores: 1
    coreLimit: 1200m
    memory: 512m
`
