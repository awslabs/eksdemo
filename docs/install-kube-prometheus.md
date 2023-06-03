# Install Kube Prometheus Stack using Amazon Managed Prometheus

The Kube Prometheus Stack uses Prometheus Operator to monitor Kubernetes. `eksdemo` supports 2 different install options:

1. `kube-prometheus-stack` — Kube Prometheus Stack
2. `kube-prometheus-stack-amp` — Kube Prometheus Stack using Amazon Managed Prometheus

There is also a Dashboard for Karpenter available:
* `kube-prometheus-karpenter-dashboards` — Karpenter Dashboards and ServiceMonitor

This tutorial walks through the installation of the Kube Prometheus Stack using Amazon Managed Prometheus. `eksdemo` will automate the following steps:
* Create an AMP workspace with alias `<cluster-name>-kube-prometheus` if one doesn't already exist.
* Create the Prometheus Server IAM Role (IRSA) to write metrics to AMP
* Create the Grafana IAM Role (IRSA) to read metrics from AMP
* Install Kube Prometheus Helm Chart with AMP configuration

The Kube Prometheus Stack includes:
* Grafana with over a dozen pre-configured dashboards
* kube-state-metrics
* prometheus-node-exporter

1. [Prerequisites](#prerequisites)
2. [Install Kube Prometheus Stack using Amazon Managed Prometheus](#install-kube-prometheus-stack-using-amazon-managed-prometheus-1)
3. [(Optional) Inspect Kube Prometheus Stack IAM Roles](#optional-inspect-kube-prometheus-stack-iam-roles)
4. [(Optional) Inspect Kube Prometheus Stack Amazon Managed Prometheus Workspace](#optional-inspect-kube-prometheus-stack-amazon-managed-prometheus-workspace)
5. [(Optional) Alternate Installation Options](#optional-alternate-installation-options)
6. [(Optional) Install Karpenter Dashboards](#optional-install-karpenter-dashboards)

## Prerequisites

Click on the name of the prerequisites in the list below for a link to detailed instructions.

* [Amazon EKS Cluster with an IAM OIDC provider configured](/docs/create-cluster.md)
* [Route 53 Public Hosted Zone](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingHostedZone.html) — Grafana will be available over the Internet using a fully qualified domain name (FQDN). This tutorial will use **`<example.com>`** and you will need to replace it with your hosted zone.
* [AWS Load Balancer Controller](/docs/install-awslb.md) — Includes `Ingress` resource to expose Grafana via an ALB.
* [ExternalDNS](/docs/install-edns.md) — ExternalDNS will add a DNS record for Grafana to your Route 53 Hosted Zone.
* [ACM Certificate](/docs/create-acm-cert.md) — Certificate for Grafana. Can be a wildcard certificate.

If you want to deploy with no prerequisites, see the [(Optional) Alternate Installation Options](#optional-alternate-installation-options)

## Install Kube Prometheus Stack using Amazon Managed Prometheus

In this section we walk through the process of installing the Kube Prometheus Stack using Amazon Managed Prometheus. The command for performing the installation is: `eksdemo install kube-prometheus-stack-amp -c <cluster-name>`

Let’s expore the command and it’s options by using the `-h` help shorthand flag.

```
» eksdemo install kube-prometheus-stack-amp -h
Install kube-prometheus-stack-amp

Usage:
  eksdemo install kube-prometheus-stack-amp [flags]

Aliases:
  kube-prometheus-stack-amp, kube-prometheus-amp

Flags:
      --chart-version string   chart version (default "46.6.0")
  -c, --cluster string         cluster to install application (required)
      --dry-run                don't install, just print out all installation steps
  -P, --grafana-pass string    grafana admin password (required)
  -h, --help                   help for kube-prometheus-stack-amp
      --ingress-class string   name of IngressClass (default "alb")
  -I, --ingress-host string    hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -n, --namespace string       namespace to install (default "monitoring")
  -X, --nginx-pass string      basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                    use NLB instead of CLB (when not using Ingress)
      --set strings            set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --target-type string     target type when deploying NLB or ALB Ingress (default "ip")
      --use-previous           use previous working chart/app versions ("34.10.0"/"v0.55.0")
  -v, --version string         application version (default "v0.65.1")
```

The Kube Prometheus Stack specific flag is:
* `--grafana-pass` — This flag sets the Grafana admin password

The `eksdemo` install of Kube Prometheus Stack using Amazon Managed Prometheus has made the following modifications to the Helm chart defaults:
* Prometheus Alert Manager is disabled
* Grafana datasource is configured to use AMP
* Grafana is configured to use AWS Sigv4 authentication
* Monitoring of etcd, kube-scheduler and kube-controller-manager is disabled
* Prometheus is configured to remote-write data to AMP
* The Prometheus remote-write dashboard is enabled
* The `cleanPrometheusOperatorObjectNames` option is used for shorter names for Prometheus Operator resources

Let's explore the dry run output with the `--dry-run` flag. The syntax for the command with all the options is: **`eksdemo install kube-prometheus-stack-amp -c <cluster-name> --ingress-host kprom.<example.com> --grafana-pass <your-password> --dry-run`**.

Please be sure to:
* Replace `<cluster-name>` with the name of your EKS cluster
* Replace `<example.com>` with your Route 53 hosted zone
* Replace `<your-password>` with a secure password

```
» eksdemo install kube-prometheus-stack-amp -c <cluster-name> --ingress-host kprom.<example.com> --grafana-pass <your-password> --dry-run
Creating 3 dependencies for kube-prometheus-stack-amp
Creating dependency: prometheus-amp-irsa

Eksctl Resource Manager Dry Run:
eksctl create iamserviceaccount -f - --approve
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: blue
  region: us-west-2

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
      name: prometheus-prometheus
      namespace: monitoring
    roleName: eksdemo.blue.monitoring.prometheus-prometheus
    roleOnly: true
    attachPolicy:
      Version: "2012-10-17"
      Statement:
      - Effect: Allow
        Action:
        - aps:RemoteWrite
        - aps:GetSeries
        - aps:GetLabels
        - aps:GetMetricMetadata
        Resource: "*"

Creating dependency: grafana-amp-irsa

Eksctl Resource Manager Dry Run:
eksctl create iamserviceaccount -f - --approve
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: blue
  region: us-west-2

iam:
  withOIDC: true
  serviceAccounts:
  - metadata:
      name: prometheus-grafana
      namespace: monitoring
    roleName: eksdemo.blue.monitoring.prometheus-grafana
    roleOnly: true
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess

Creating dependency: amazon-managed-prometheus-workspace

AMP Resource Manager Dry Run:
Amazon Managed Service for Prometheus API Call "CreateWorkspace" with request parameters:
alias: "blue-kube-prometheus"

Helm Installer Dry Run:
+---------------------+----------------------------------------------------+
| Application Version | v0.65.1                                            |
| Chart Version       | 46.6.0                                             |
| Chart Repository    | https://prometheus-community.github.io/helm-charts |
| Chart Name          | kube-prometheus-stack                              |
| Release Name        | kube-prometheus-stack-amp                          |
| Namespace           | monitoring                                         |
| Wait                | false                                              |
+---------------------+----------------------------------------------------+
Set Values: []
Values File:
---
fullnameOverride: prometheus
defaultRules:
  rules:
    alertmanager: false
alertmanager:
  enabled: false
grafana:
  fullnameOverride: grafana
  serviceAccount:
    name: prometheus-grafana
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.monitoring.prometheus-grafana
  ingress:
    enabled: true
    ingressClassName: alb
    annotations:
      alb.ingress.kubernetes.io/listen-ports: '[{"HTTP": 80}, {"HTTPS":443}]'
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/ssl-redirect: '443'
      alb.ingress.kubernetes.io/target-type: ip
    hosts:
    - kprom.example.com
    tls:
    - hosts:
      - kprom.example.com
  sidecar:
    datasources:
      defaultDatasourceEnabled: false
  additionalDataSources:
  - name: Amazon Managed Service for Prometheus
    type: prometheus
    url: <-amp_endpoint_url_will_go_here->
    access: proxy
    isDefault: true
    jsonData:
      sigV4Auth: true
      sigV4AuthType: default
      sigV4Region: us-west-2
    timeInterval: 30s
  service:
    annotations:
      {}
    type: ClusterIP
  # Temporary fix for issue: https://github.com/prometheus-community/helm-charts/issues/1867
  # Note the "serviceMonitorSelectorNilUsesHelmValues: false" below also resolves the issue
  serviceMonitor:
    labels:
      release: kube-prometheus-stack-amp
  adminPassword: password
  grafana.ini:
    auth:
      sigv4_auth_enabled: true
kubeControllerManager:
  enabled: false
kubeEtcd:
  enabled: false
kubeScheduler:
  enabled: false
kube-state-metrics:
  fullnameOverride: kube-state-metrics
prometheus-node-exporter:
  fullnameOverride: node-exporter
prometheusOperator:
  image:
    tag: v0.65.1
prometheus:
  serviceAccount:
    name: prometheus-prometheus
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.monitoring.prometheus-prometheus
  prometheusSpec:
    # selects ServiceMonitors without the "release: kube-prometheus-stack-amp" label
    serviceMonitorSelectorNilUsesHelmValues: false
    remoteWrite:
    - url: <-amp_endpoint_url_will_go_here->api/v1/remote_write
      sigv4:
        region: us-west-2
      queueConfig:
        maxSamplesPerSend: 1000
        maxShards: 200
        capacity: 2500
    remoteWriteDashboards: true
cleanPrometheusOperatorObjectNames: true

Helm Installer Post Render Kustomize Dry Run:
---
resources:
- manifest.yaml
patches:
# Delete the AlertManager dashboard as it's disabled
- patch: |-
    $patch: delete
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: needed-but-not-used
  target:
    name: ".*alertmanager-overview$"
```

In the `--dry-run` output above, you will notice a placeholder:
* <-amp_endpoint_url_will_go_here->
This will be replaced with the actual AMP Workspace Endpoint URL after the Workspace is created.

Let's proceed with installing the Kube Prometheus Stack using Amazon Managed Prometheus. Please be sure to:
* Replace `<cluster-name>` with the name of your EKS cluster
* Replace `<example.com>` with your Route 53 hosted zone
* Replace `<your-password>` with a secure password

```
» eksdemo install kube-prometheus-stack-amp -c <cluster-name> --ingress-host kprom.<example.com> --grafana-pass <your-password>
Creating 3 dependencies for kube-prometheus-stack-amp
Creating dependency: prometheus-amp-irsa
<snip>
Creating dependency: grafana-amp-irsa
<snip>
Creating dependency: amazon-managed-prometheus-workspace
Creating AMP Workspace Alias: blue-kube-prometheus...done
Created AMP Workspace Id: ws-a29b716d-7f17-4689-a0b4-26073c6790a0
Downloading Chart: https://github.com/prometheus-community/helm-charts/releases/download/kube-prometheus-stack-46.6.0/kube-prometheus-stack-46.6.0.tgz
Helm installing...
<snip>
Using chart version "46.6.0", installed "kube-prometheus-stack-amp" version "v0.65.1" in namespace "monitoring"
NOTES:
kube-prometheus-stack has been installed. Check its status by running:
  kubectl --namespace monitoring get pods -l "release=kube-prometheus-stack-amp"

Visit https://github.com/prometheus-operator/kube-prometheus for instructions on how to create & configure Alertmanager and Prometheus instances using the Operator.
```

Open your web browser and enter `https://kprom.<example.com>` (**replace `<example.com>` with your Hosted Zone**) to load Grafana. When asked to sign-in, enter `admin` for username and the password you used with the `--grafana-pass` flag for password. 

![Kube Prometheus Dashboard](/docs/images/kube-prometheus.png "Kube Prometheus Dashboard")

NOTE: It’s possible you may have to wait for DNS to propagate. The time depends on your local ISP and operating system. If you get a DNS resolution error, you can wait and try again later. Or if you’d like to troubleshoot a bit further, A2 Hosting has a Knowledge base article [How to test DNS with dig and nslookup](https://www.a2hosting.com/kb/getting-started-guide/internet-and-networking/troubleshooting-dns-with-dig-and-nslookup).

## (Optional) Inspect Kube Prometheus Stack IAM Roles

The Kube Prometheus Stack using Amazon Managed Prometheus install creates 2 IAM Roles:
* Prometheus Server IAM Role (IRSA)
* Grafana IAM Role (IRSA)

Use the `eksdemo get iam-role` command and the `--search` flag to find the roles.

```
» eksdemo get iam-role --search prometheus
+------------+-----------------------------------------------+
|    Age     |                     Role                      |
+------------+-----------------------------------------------+
| 14 minutes | eksdemo.blue.monitoring.prometheus-grafana    |
| 15 minutes | eksdemo.blue.monitoring.prometheus-prometheus |
+------------+-----------------------------------------------+
```

`eksdemo` uses a specific naming convention for IRSA roles: `eksdemo.<cluster-name>.<namespace>.<serviceaccount-name>`. To view the permissions assigned to the role, use the `eksdemo get iam-policy` command and the `--role` command which lists only the policies assigned to the role. 

Let's view the Prometheus Server IAM Role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role eksdemo.<cluster-name>.monitoring.prometheus-prometheus 
+------------------------------------------------------------------------------+--------+-------------+
|                                     Name                                     |  Type  | Description |
+------------------------------------------------------------------------------+--------+-------------+
| eksctl-blue-addon-iamserviceaccount-monitoring-prometheus-prometheus-Policy1 | Inline |             |
+------------------------------------------------------------------------------+--------+-------------+
```

To view the details of the policy, including the policy document, you can use the `--output` flag or `-o` shorthand flag to output the raw AWS API responses in either JSON or YAML. In the example we'll use YAML.

```
» eksdemo get iam-policy --role  eksdemo.<cluster-name>.kubecost.kubecost-cost-analyzer -o yaml
InlinePolicies:
- Name: eksctl-blue-addon-iamserviceaccount-monitoring-prometheus-prometheus-Policy1
  PolicyDocument: |-
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Action": [
                    "aps:RemoteWrite",
                    "aps:GetSeries",
                    "aps:GetLabels",
                    "aps:GetMetricMetadata"
                ],
                "Resource": "*",
                "Effect": "Allow"
            }
        ]
    }
ManagedPolicies: []
```

You can run the same commands to view the Grafana IAM role details.

## (Optional) Inspect Kube Prometheus Stack Amazon Managed Prometheus Workspace

The Kube Prometheus Stack install creates an Amazon Managed Prometheus (AMP) workspace to store the Kubernetes cluster performance metrics. Use the `eksdemo get amp-workspace` command to inspect the AMP workspace that was created.

```
» eksdemo get amp-workspace
+------------+--------+----------------------+-----------------------------------------+
|    Age     | Status |        Alias         |              Workspace Id               |
+------------+--------+----------------------+-----------------------------------------+
| 31 minutes | ACTIVE | blue-kube-prometheus | ws-a29b716d-7f17-4689-a0b4-26073c6790a0 |
+------------+--------+----------------------+-----------------------------------------+
```

To view the raw output of the AWS API response use the `-o yaml` output option. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get amp-workspace <cluster-name>-kube-prometheus -o yaml
- Workspace:
    Alias: blue-kube-prometheus
    Arn: arn:aws:aps:us-west-2:123456789012:workspace/ws-a29b716d-7f17-4689-a0b4-26073c6790a0
    CreatedAt: "2023-03-06T02:32:37.164Z"
    PrometheusEndpoint: https://aps-workspaces.us-west-2.amazonaws.com/workspaces/ws-a29b716d-7f17-4689-a0b4-26073c6790a0/
    Status:
      StatusCode: ACTIVE
    Tags: {}
    WorkspaceId: ws-a29b716d-7f17-4689-a0b4-26073c6790a0
  WorkspaceLogging: null
```

## (Optional) Alternate Installation Options

`eksdemo` makes it easy to make changes to the Kube Prometheus Stack install with 2 different application installs and different command line flags. The following are just a few examples of different install options:

* Install Kube Prometheus Stack using Amazon Managed Prometheus without a domain with a CLB
    * `eksdemo install kube-prometheus-stack-amp -c <cluster-name> -P <your-password>`
* Install Kube Prometheus Stack using Amazon Managed Prometheus without a domain with a NLB
    * `eksdemo install kube-prometheus-stack-amp -c <cluster-name> -P <your-password> --nlb`
* Install Kube Prometheus Stack with Nginx Ingress
    * `eksdemo install kube-prometheus-stack -c <cluster-name> -P <your-password> -I kprom.<example.com> --ingress-class=nginx`

## (Optional) Install Karpenter Dashboards

`eksdemo` make it easy to add [Karpenter Dashboards](https://karpenter.sh/docs/getting-started/getting-started-with-eksctl/#add-optional-monitoring-with-grafana) to the Kube Prometheus Stack install with a single command that will:
* Create the Karpenter `ServiceMonitor` which instructs Prometheus to scrape Karpenter metrics
* Add the Karpenter Dashboards to the Kube Prometheus Grafana instance:
    * Karpenter Capacity Dashboard
    * Karpenter Performance Dashboard

The command is `eksdemo install kube-prometheus-karpenter-dashboards -c <cluster-name>`. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install kube-prometheus-karpenter-dashboards -c blue
Helm installing...
2023/03/05 20:25:30 creating 1 resource(s)
2023/03/05 20:25:30 creating 3 resource(s)
Using chart version "n/a", installed "kube-prometheus-karpenter-dashboards" version "n/a" in namespace "monitoring"
```

![Karpenter Performance Dashboard](/docs/images/karpenter-perf-dashboard.png "Karpenter Performance Dashboard")
