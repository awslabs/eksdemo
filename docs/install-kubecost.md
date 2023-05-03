# Install EKS optimized Kubecost using Amazon Managed Prometheus

Kubecost helps customers monitor and manage Kubernetes spend. `eksdemo` supports 3 different install options for Kubecost:

1. `kubecost-vendor` — Vendor distribution of Kubecost
2. `kubecost-eks` — EKS optimized bundle of Kubecost
3. `kubecost-eks-amp` — EKS optimized Kubecost using Amazon Managed Prometheus

This tutorial walks through the 3rd option, the installation of the EKS optimized Kubecost using Amazon Managed Prometheus (AMP). `eksdemo` makes it very easy and with a single command will automate the following steps:
* Create an AMP workspace with alias `<cluster-name>-kubecost` if one doesn't already exist.
* Create the Kubecost Cost Analyzer IAM Role (IRSA)
* Create the Kubecost Prometheus Server IAM Role (IRSA)
* Install Kubecost configured to use the AMP workspace
* Secure the Kubecost dashboard behind a password (HTTP basic authentication)

1. [Prerequisites](#prerequisites)
2. [Install EKS optimized Kubecost using Amazon Managed Prometheus](#install-eks-optimized-kubecost-using-amazon-managed-prometheus-1)
3. [(Optional) Inspect Kubecost IAM Roles](#optional-inspect-kubecost-iam-roles)
4. [(Optional) Inspect Kubecost Amazon Managed Prometheus Workspace](#optional-inspect-kubecost-amazon-managed-prometheus-workspace)
5. [(Optional) Alternate Installation Options](#optional-alternate-installation-options)

## Prerequisites

Click on the name of the prerequisites in the list below for a link to detailed instructions.

* [Amazon EKS Cluster with an IAM OIDC provider configured](/docs/create-cluster.md)
* [EBS CSI Driver](/docs/install-ebs-csi-driver.md) — The EBS CSI Friver is required if your cluster is running EKS versions 1.23 or higher to provision the EBS volume for Promtheus.
* [Route 53 Public Hosted Zone](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingHostedZone.html) — Kubecost will be available over the Internet using a fully qualified domain name (FQDN). This tutorial will use **`<example.com>`** and you will need to replace it with your hosted zone.
* [ExternalDNS](/docs/install-edns.md) — ExternalDNS will add a DNS record for Kubecost to your Route 53 Hosted Zone.
* [Ingress NGINX](/docs/install-ingress-nginx.md) — Ingress NGINX will serve as the Ingress controller for Kubecost.
* [cert-manager](/docs/install-cert-manager.md) — cert-manager will create a publicly trusted certificate for Ingress NGINX.

Ingress NGINX was choosen as the Ingress controller because it offers simple user/password authentication. This allows you to secure the Kubecost application behind a password.

If you don't want to use Ingress NGINX or want to deploy without a Route 53 Pulic Hosted Zone, see the [(Optional) Alternate Installation Options](#optional-alternate-installation-options)

## Install EKS optimized Kubecost using Amazon Managed Prometheus

In this section we walk through the process of installing the EKS optimized bundle of Kubecost using Amazon Managed Prometheus. The command for performing the installation is: `eksdemo install kubecost-eks-amp -c <cluster-name>`

Let’s expore the command and it’s options by using the `-h` help shorthand flag.

```
» eksdemo install kubecost-eks-amp -h
Install kubecost-eks-amp

Usage:
  eksdemo install kubecost-eks-amp [flags]

Flags:
      --chart-version string     chart version (default "1.100.0")
  -c, --cluster string           cluster to install application (required)
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for kubecost-eks-amp
      --ingress-class string     name of IngressClass (default "alb")
  -I, --ingress-host string      hostname for Ingress with TLS (default is Service of type LoadBalancer)
  -n, --namespace string         namespace to install (default "kubecost")
  -X, --nginx-pass string        basic auth password for admin user (only valid with --ingress-class=nginx)
      --nlb                      use NLB instead of CLB (when not using Ingress)
      --no-prometheus            don't install prometheus
      --node-exporter            install prometheus node exporter (not installed by default)
      --service-account string   service account name (default "kubecost-cost-analyzer")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --target-type string       target type when deploying NLB or ALB Ingress (default "ip")
      --use-previous             use previous working chart/app versions ("1.97.0"/"1.97.0")
  -v, --version string           application version (default "1.100.0")
```

The Kubecost specific flags are:
* `--no-prometheus` — This flag will configure Kubecost to install without the local Prometheus server. Kubecost is still configured to use Amazon Managed Prometheus. You will need your own local Prometheus or ADOT setup with the required scrape configs.
* `--node-exporter` — The EKS optmized bundle of Kubecost is configured with Prometheus Node Exporter disabled. This is to prevent a conflict if Node Exporter is already deployed in the cluster. Use this flag if Node Exporter is not installed in your cluster to include Prometheus Node Exporter as part of the Kubecost install.

There are three flags used for the install that aren't Kubecost specific:
* `--ingress-class` — `eksdemo` defaults to using the `alb` Ingress class which depends on the AWS Load Balancer Controller to deploy an Application Load Balancer (ALB). Since we are using Ingress Nginx, we will use this flag to specify the `nginx` Ingress class.
* `--ingress-host` — This flag is used to specify the fully qualified domain name for the application. It's used as the host component in the Ingress definition. This tutorial will use **`kubecost.<example.com>`**
* `--nginx-pass` — This flag configures HTTP basic authentication and used to specify the password for the `admin` user.

The `eksdemo` install of EKS optimized Kubecost using Amazon Managed Prometheus is based on the Amazon Managed Prometheus documentation section [Integrating with Amazon EKS cost monitoring](https://docs.aws.amazon.com/prometheus/latest/userguide/integrating-kubecost.html). The documentation uses two custom values files in the [Kubecost Helm Chart Repository](https://github.com/kubecost/cost-analyzer-helm-chart/tree/develop/cost-analyzer):
* [values-eks-cost-monitoring.yaml](https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-eks-cost-monitoring.yaml)
* [values-amp.yaml](https://github.com/kubecost/cost-analyzer-helm-chart/blob/develop/cost-analyzer/values-amp.yaml)

Let's explore the dry run output with the `--dry-run` flag. The syntax for the command with all the options is: **`eksdemo install kubecost-eks-amp -c <cluster-name> --node-exporter --ingress-host kubecost.<example.com> --ingress-class=nginx --nginx-pass <your-password> --dry-run`**.

Please be sure to:
* Replace `<cluster-name>` with the name of your EKS cluster
* Replace `<example.com>` with your Route 53 hosted zone
* Replace `<your-password>` with a secure password

```
» eksdemo install kubecost-eks-amp -c <cluster-name> --node-exporter --ingress-host kubecost.<example.com> --ingress-class=nginx --nginx-pass <your-password> --dry-run
Creating 3 dependencies for kubecost-eks-amp

Creating dependency: kubecost-cost-analyzer-irsa

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
      name: kubecost-cost-analyzer
      namespace: kubecost
    roleName: eksdemo.blue.kubecost.kubecost-cost-analyzer
    roleOnly: true
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess
    - arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess


Creating dependency: kubecost-prometheus-server-irsa

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
      name: kubecost-prometheus-server
      namespace: kubecost
    roleName: eksdemo.blue.kubecost.kubecost-prometheus-server
    roleOnly: true
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess
    - arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess


Creating dependency: kubecost-amazon-managed-prometheus

AMP Resource Manager Dry Run:
Amazon Managed Service for Prometheus API Call "CreateWorkspace" with request parameters:
alias: "blue-kubecost"

Helm Installer Dry Run:
+---------------------+---------------------------------------------+
| Application Version | 1.100.0                                     |
| Chart Version       | 1.100.0                                     |
| Chart Repository    | oci://public.ecr.aws/kubecost/cost-analyzer |
| Chart Name          | cost-analyzer                               |
| Release Name        | kubecost-eks-amp                            |
| Namespace           | kubecost                                    |
| Wait                | false                                       |
+---------------------+---------------------------------------------+
Set Values: []
Values File:
---
fullnameOverride: kubecost-cost-analyzer
global:
  amp:
    enabled: true
    prometheusServerEndpoint: http://localhost:8005/workspaces/<-amp_id_will_go_here->
    remoteWriteService: <-amp_endpoint_url_will_go_here->api/v1/remote_write
    sigv4:
      region: us-west-2
  grafana:
    # If false, Grafana will not be installed
    enabled: false
    # If true, the kubecost frontend will route to your grafana through its service endpoint
    proxy: false
  prometheus:
    enabled: true
sigV4Proxy:
  region: us-west-2
  host: aps-workspaces.us-west-2.amazonaws.com
podSecurityPolicy:
  enabled: false
imageVersion: prod-1.100.0
kubecostFrontend:
  image: public.ecr.aws/kubecost/frontend
kubecostModel:
  image: public.ecr.aws/kubecost/cost-model
  # The total number of days the ETL storage will build
  etlStoreDurationDays: 120
serviceAccount:
  name: kubecost-cost-analyzer
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.kubecost.kubecost-cost-analyzer
ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: basic-auth
    nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"
  pathType: Prefix
  hosts:
    - kubecost.example.com
  tls:
  - hosts:
    - kubecost.example.com
    secretName: cost-analyzer-tls
service:
  type: ClusterIP
  annotations:
    {}
prometheus:
  serviceAccounts:
    server:
      annotations:
        eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.kubecost.kubecost-prometheus-server
  kube-state-metrics:
    fullnameOverride: kubecost-kube-state-metrics
  server:
    fullnameOverride: kubecost-prometheus-server
    image:
      repository: public.ecr.aws/kubecost/prometheus
      tag: v2.35.0
    global:
      # overrides kubecost default of 60s, sets to prom default of 10s
      scrape_timeout: 10s
      external_labels:
        cluster_id: blue # Each cluster should have a unique ID
  configmapReload:
    prometheus:
      enabled: false
      image:
        repository: public.ecr.aws/bitnami/configmap-reload
        tag: 0.7.1
  nodeExporter:
    enabled: true
    fullnameOverride: kubecost-prometheus-node-exporter
reporting:
  productAnalytics: false
kubecostProductConfigs:
  clusterName: blue

Creating 1 post-install resources for kubecost-eks-amp
Creating post-install resource: kubecost-eks-amp-basic-auth-secret

Kubernetes Resource Manager Dry Run:
---
apiVersion: v1
kind: Secret
metadata:
  name: basic-auth
  namespace: kubecost
type: Opaque
data:
  auth: YWRtaW46JDJhJDEwJFl1ZFhjcTVqdzZXNXNEWE5lTi8ycE8xdXdNZ05lcy5VNmZaZm5mV0pXT2JoanJ4SHI0QW1D
```

In the `--dry-run` output above, you will notice two placeholders:
* <-amp_id_will_go_here->
* <-amp_endpoint_url_will_go_here->
These will be replaced with the actual AMP Workspace Id and Endpoint URL after the Workspace is created.

There is also a Kubernetes Secret that is created post-install called `basic-auth` that configures the password for Ingress NGINX to setup HTTP Basic Authentication.

Let's proceed with installing Kubecost. Please be sure to:
* Replace `<cluster-name>` with the name of your EKS cluster
* Replace `<example.com>` with your Route 53 hosted zone
* Replace `<your-password>` with a secure password

```
» eksdemo install kubecost-eks-amp -c <cluster-name> --node-exporter --ingress-host kubecost.<example.com> --ingress-class=nginx --nginx-pass <your-password>
Creating 3 dependencies for kubecost-eks-amp

Creating dependency: kubecost-cost-analyzer-irsa
<snip>

Creating dependency: kubecost-prometheus-server-irsa
<snip>

Creating dependency: kubecost-amazon-managed-prometheus
Creating AMP Workspace Alias: blue-kubecost...done
Created AMP Workspace Id: ws-1c43b388-1ff6-4fdd-8ba9-2493fe287dad
Downloading Chart: oci://public.ecr.aws/kubecost/cost-analyzer:1.100.0
Helm installing...
2023/02/11 19:42:13 creating 1 resource(s)
2023/02/11 19:42:14 creating 40 resource(s)
Using chart version "1.100.0", installed "kubecost-eks-amp" version "1.100.0" in namespace "kubecost"
NOTES:
--------------------------------------------------Kubecost has been successfully installed.

WARNING: ON EKS v1.23+ INSTALLATION OF EBS-CSI DRIVER IS REQUIRED TO MANAGE PERSISTENT VOLUMES. LEARN MORE HERE: https://docs.kubecost.com/install-and-configure/install/provider-installations/aws-eks-cost-monitoring#prerequisites

Please allow 5-10 minutes for Kubecost to gather metrics.

If you have configured cloud-integrations, it can take up to 48 hours for cost reconciliation to occur.

When using Durable storage (Enterprise Edition), please allow up to 4 hours for data to be collected and the UI to be healthy.

When pods are Ready, you can enable port-forwarding with the following command:

    kubectl port-forward --namespace kubecost deployment/kubecost-cost-analyzer 9090

Next, navigate to http://localhost:9090 in a web browser.

Having installation issues? View our Troubleshooting Guide at http://docs.kubecost.com/troubleshoot-install
Creating 1 post-install resources for kubecost-eks-amp
Creating post-install resource: kubecost-eks-amp-basic-auth-secret
Creating Secret "basic-auth" in namespace "kubecost"
```

Open your web browser and enter `https://kubecost.<example.com>` (**replace `<example.com>` with your Hosted Zone**) to load the Kubecost dashboard. When asked to sign-in, enter `admin` for username and the password you used with the `--nginx-pass` flag for password. 

![Kubecost Screenshot](/docs/images/kubecost-screenshot.jpg "Kubecost Screenshot")

NOTE: It’s possible you may have to wait for DNS to propagate. The time depends on your local ISP and operating system. If you get a DNS resolution error, you can wait and try again later. Or if you’d like to troubleshoot a bit further, A2 Hosting has a Knowledge base article [How to test DNS with dig and nslookup](https://www.a2hosting.com/kb/getting-started-guide/internet-and-networking/troubleshooting-dns-with-dig-and-nslookup).

## (Optional) Inspect Kubecost IAM Roles 

The Kubecost install creates 2 IAM Roles:
* Kubecost Cost Analyzer IAM Role (IRSA)
* Kubecost Prometheus Server IAM Role (IRSA)

Use the `eksdemo get iam-role` command and the `--search` flag to find the roles.

```
» eksdemo get iam-role --search kubecost
+------------+--------------------------------------------------+
|    Age     |                       Role                       |
+------------+--------------------------------------------------+
| 45 minutes | eksdemo.blue.kubecost.kubecost-cost-analyzer     |
| 44 minutes | eksdemo.blue.kubecost.kubecost-prometheus-server |
+------------+--------------------------------------------------+
```

`eksdemo` uses a specific naming convention for IRSA roles: `eksdemo.<cluster-name>.<namespace>.<serviceaccount-name>`. To view the permissions assigned to the role, use the `eksdemo get iam-policy` command and the `--role` command which lists only the policies assigned to the role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role eksdemo.<cluster-name>.kubecost.kubecost-cost-analyzer
+-----------------------------------+---------+--------------------------------+
|               Name                |  Type   |          Description           |
+-----------------------------------+---------+--------------------------------+
| AmazonPrometheusRemoteWriteAccess | AWS Mgd | Grants write only access       |
|                                   |         | to AWS Managed Prometheus      |
|                                   |         | workspaces                     |
+-----------------------------------+---------+--------------------------------+
| AmazonPrometheusQueryAccess       | AWS Mgd | Grants access to run queries   |
|                                   |         | against AWS Managed Prometheus |
|                                   |         | resources                      |
+-----------------------------------+---------+--------------------------------+
```

To view the details of the policy, including the policy document, you can use the `--output` flag or `-o` shorthand flag to output the raw AWS API responses in either JSON or YAML. In the example we'll use YAML.

```
» eksdemo get iam-policy --role  eksdemo.<cluster-name>.kubecost.kubecost-cost-analyzer -o yaml
InlinePolicies: []
ManagedPolicies:
- Policy:
    Arn: arn:aws:iam::aws:policy/AmazonPrometheusRemoteWriteAccess
    AttachmentCount: 2
    CreateDate: "2020-12-19T01:04:32Z"
    DefaultVersionId: v1
    Description: Grants write only access to AWS Managed Prometheus workspaces
    IsAttachable: true
    Path: /
    PermissionsBoundaryUsageCount: 0
    PolicyId: ABCDEFGHIJKLMNOPQ1234
    PolicyName: AmazonPrometheusRemoteWriteAccess
    Tags: []
    UpdateDate: "2020-12-19T01:04:32Z"
  PolicyVersion:
    CreateDate: "2020-12-19T01:04:32Z"
    Document: |-
      {
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Action": [
                      "aps:RemoteWrite"
                  ],
                  "Effect": "Allow",
                  "Resource": "*"
              }
          ]
      }
    IsDefaultVersion: true
    VersionId: v1
- Policy:
    Arn: arn:aws:iam::aws:policy/AmazonPrometheusQueryAccess
    AttachmentCount: 2
    CreateDate: "2020-12-19T01:02:58Z"
    DefaultVersionId: v1
    Description: Grants access to run queries against AWS Managed Prometheus resources
    IsAttachable: true
    Path: /
    PermissionsBoundaryUsageCount: 0
    PolicyId: ABCDEFGHIJKLMNOPQ1234
    PolicyName: AmazonPrometheusQueryAccess
    Tags: []
    UpdateDate: "2020-12-19T01:02:58Z"
  PolicyVersion:
    CreateDate: "2020-12-19T01:02:58Z"
    Document: |-
      {
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Action": [
                      "aps:RemoteWrite"
                  ],
                  "Effect": "Allow",
                  "Resource": "*"
              }
          ]
      }{
          "Version": "2012-10-17",
          "Statement": [
              {
                  "Action": [
                      "aps:GetLabels",
                      "aps:GetMetricMetadata",
                      "aps:GetSeries",
                      "aps:QueryMetrics"
                  ],
                  "Effect": "Allow",
                  "Resource": "*"
              }
          ]
      }
    IsDefaultVersion: true
    VersionId: v1
```

Feel free to run the same commands to view the kubecost-prometheus-server IAM role details.

## (Optional) Inspect Kubecost Amazon Managed Prometheus Workspace

The Kubecost install creates an Amazon Managed Prometheus (AMP) workspace to store the Prometheus metrics used to calculate the usage metrics and cost. Use the `eksdemo get amp-workspace` command to inspect the AMP workspace that was created.

```
» eksdemo get amp-workspace
+------------+--------+---------------+-----------------------------------------+
|    Age     | Status |     Alias     |              Workspace Id               |
+------------+--------+---------------+-----------------------------------------+
| 54 minutes | ACTIVE | blue-kubecost | ws-1c43b388-1ff6-4fdd-8ba9-2493fe287dad |
+------------+--------+---------------+-----------------------------------------+
```

To view the raw output of the AWS API response use the `-o yaml` output option. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get amp-workspace <cluster-name>-kubecost -o yaml
- Workspace:
    Alias: blue-kubecost
    Arn: arn:aws:aps:us-west-2:123456789012:workspace/ws-1c43b388-1ff6-4fdd-8ba9-2493fe287dad
    CreatedAt: "2023-02-12T02:41:53.561Z"
    PrometheusEndpoint: https://aps-workspaces.us-west-2.amazonaws.com/workspaces/ws-1c43b388-1ff6-4fdd-8ba9-2493fe287dad/
    Status:
      StatusCode: ACTIVE
    Tags: {}
    WorkspaceId: ws-1c43b388-1ff6-4fdd-8ba9-2493fe287dad
  WorkspaceLogging: null
```

## (Optional) Alternate Installation Options

`eksdemo` makes it easy to make changes to the Kubecost install with 3 different application installs and different command line flags. The following are just a few examples of different install options:

* Install EKS optimized Kubecost using Amazon Managed Prometheus without a domain with a CLB
    * `eksdemo install kubecost-eks-amp -c <cluster-name> --node-exporter`
* Install EKS optimized Kubecost using Amazon Managed Prometheus without a domain with a NLB
    * `eksdemo install kubecost-eks-amp -c <cluster-name> --node-exporter --nlb`
* Install EKS optimized Kubecost with an ALB
    * `eksdemo install kubecost-eks -c <cluster-name> --node-exporter -I kubecost.<example.com>`
* Install the vendor version of Kubecost with an ALB
    * `eksdemo install kubecost-vendor -c <cluster-name> -I kubecost.<example.com>`

Just a few considations:
* If you are exposing Kubecost with an ALB, you will need:
    * [AWS Load Balancer Controller](/docs/install-awslb.md)
    * [ExternalDNS](/docs/install-edns.md)
    * [ACM Certificate](/docs/create-acm-cert.md)
* If you already have Prometheus Node Exporter installed in your cluster, remove the `--node-exporter` flag
