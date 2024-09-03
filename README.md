# `eksdemo` - kubectl-like CLI for Amazon EKS

[![Go Report Card](https://goreportcard.com/badge/github.com/awslabs/eksdemo)](https://goreportcard.com/report/github.com/awslabs/eksdemo)

The easy button for learning, testing, and demoing Amazon EKS:
* Install complex applications and dependencies with a single command
* Extensive application catalog with over 70 options (includes CNCF, open source and related projects)
* Customize application installs easily with simple command line flags
* Query and search AWS resources with over 60 kubectl-like get commands

> Note: `eksdemo` is in beta and is intended for demo and test environments only.

## Table of Contents
- [Why `eksdemo`?](#why-eksdemo)
  - [No Magic](#no-magic)
  - [`eksdemo` vs EKS Blueprints](#eksdemo-vs-eks-blueprints)
- [Install `eksdemo`](#install-eksdemo)
  - [Prerequisites](#prerequisites)
  - [Install using Homebrew](#install-using-homebrew)
  - [Troubleshoot Homebrew Install](#troubleshoot-homebrew-install)
  - [Install on AWS CloudShell](#install-on-aws-cloudshell)
  - [Install on AWS Cloud9](#install-on-aws-cloud9)
  - [Install Manually](#install-manually)
  - [Set the AWS Region](#set-the-aws-region)
  - [Validate Install](#validate-install)
- [Application Catalog](#application-catalog)
- [Kubectl-like get commands](#kubectl-like-get-commands)
- [Tutorials](#tutorials)
  - [Basics](#basics)
  - [Advanced](#advanced)
- [Support \& Feedback](#support--feedback)
- [Security](#security)
- [License](#license)

## Why `eksdemo`?
While creating an EKS cluster is fairly easy thanks to [`eksctl`](https://eksctl.io/), manually installing and configuring applications on EKS is complex, time consuming and error-prone. One of the most powerful feature of `eksdemo` is its extensive application catalog. An application can be installed (including dependencies) with a single command.

For example, the command: **`eksdemo install karpenter -c <cluster-name>`** will:
1. Create the EC2 Spot Service Linked Role (if it doesn't already exist)
2. Create the Karpenter Controller IAM Role (IRSA)
3. Create the Karpenter Node IAM Role
4. Create an SQS Queue and EventBridge rules for native Spot Termination Handling
5. Add an entry to the `aws-auth` ConfigMap for the Karpenter Node IAM Role
6. Install the Karpenter Helm Chart
7. Create default Karpenter `NodePool` and `EC2NodeClass` Custom Resources

### No Magic
Application installs are:
* Transparent
    * The `--dry-run` flag prints out all the steps `eksdemo` will take to create dependencies and install the application
* Customizable
    * Each application has optional flags for common configuration options
    * The `--set` flag is available to override any settings in a Helm chart's values file 
* Managed by Helm
    * `eksdemo` embeds Helm as a library and it's used to install all applications, even those that don't have a Helm chart

### `eksdemo` vs EKS Blueprints

Both `eksdemo` and [EKS Blueprints](https://aws.amazon.com/blogs/containers/bootstrapping-clusters-with-eks-blueprints/) automate the creation of EKS clusters and install commonly used applications. Why would you use `eksdemo` for learning, testing, and demoing EKS?

| `eksdemo` | EKS Blueprints |
------------|-----------------
Use cases: learning, testing, and demoing EKS | Use cases: customers deploying to prod and non-prod environments
Kubectl-like CLI installs apps with single command | Infrastructure as Code (IaC) built on Terraform or CDK
Imperative tooling is great for iterative testing | Declarative IaC tooling is not designed for iterative testing
Used to get up and running quickly | Used to drive standards and communicate vetted architecture patterns  for utilizing EKS within customer organizations

## Install `eksdemo`

`eksdemo` is a Golang binary and releases include support for Mac, Linux and Windows running on x86 or arm64. There are a number of ways you can install:

* [Install using Homebrew](#install-using-homebrew) — Install locally for Mac and Linux users.
* [Install on AWS CloudShell](#install-on-aws-cloudshell) — Easiest way to use `eksdemo`.
* [Install on AWS Cloud9](#install-on-aws-cloud9) — Install into your Cloud9 instance.
* [Install Manually](#install-manually) — This method is required for Windows users.

### Prerequisites

1. AWS Account with Administrator access
2. Route53 Public Hosted Zone (Optional but strongly recommended)
    1. You can update the domain registration of your existing domain (using any domain registrar) to [change the name servers for the domain to use the four Route 53 name servers](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/migrate-dns-domain-inactive.html#migrate-dns-update-domain-inactive). 
    2. You can still use `eksdemo` if you don’t have a Route53 Hosted Zone. Most applications that explose an Ingress resource default to deploying a Service of type LoadBalancer if you don't use the `--ingress-host` flag and your connection to the application will be unencrypted.


### Install using Homebrew

[Homebrew](https://brew.sh/) installation method is supported for Mac and Linux. Using the Terminal, enter the following commands:

```
brew tap aws/tap
brew install eksdemo
```

### Troubleshoot Homebrew Install

Note: Depending on how you originally installed `eksctl`, you may receive the error: `eksctl is already installed from homebrew/core!`  This is because `eksdemo` uses the official Weaveworks tap `weaveworks/tap` as a dependency.

If you receive the error above, run the following commands:

```
brew uninstall eksctl
brew install eksdemo
```

### Install on AWS CloudShell

AWS CloudShell environments are mostly ephemeral and software you install is lost the next time you reconnect. There is [persistent storage available](https://docs.aws.amazon.com/cloudshell/latest/userguide/limits.html#persistent-storage-limitations) in the home directory that is retained for 120 days after the end of your last session. Use the following commands to install `eksdemo` on AWS CloudShell:

```
curl -s -L "https://github.com/weaveworks/eksctl/releases/latest/download/eksctl_Linux_amd64.tar.gz" | tar xz -C /tmp
curl -s -L "https://github.com/awslabs/eksdemo/releases/latest/download/eksdemo_Linux_x86_64.tar.gz" | tar xz -C /tmp
mkdir -p ~/.local/bin && mv /tmp/eksctl ~/.local/bin && mv /tmp/eksdemo ~/.local/bin
```

To configure Bash completion, first install the bash-completion package:

```
sudo dnf install bash-completion -y
```

Once that completes, run the following commands that will configure bash completion to work across CloudShell sessions by installing everything needed in the home directory:

```
mkdir -p ~/.bashrc.d
cp /usr/share/bash-completion/bash_completion ~/.bashrc.d/
echo '. <(eksdemo completion bash)' >> ~/.bashrc
source ~/.bashrc
```

### Install on AWS Cloud9

To install on Cloud9, you will need to install Homebrew first. In the Cloud9 terminal, run the following commands:

```
CI=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
(echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> ~/.bashrc
eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
brew install aws/tap/eksdemo
```

To configure Bash completion, run the following commands and then start a new terminal:
```
cat >> ~/.bashrc << \EOF
if type brew &>/dev/null
then
  for COMPLETION in "$(brew --prefix)/etc/bash_completion.d/"*
  do
    [[ -r "${COMPLETION}" ]] && source "${COMPLETION}"
  done
fi
EOF
```

### Install Manually

Navigate to [Releases](https://github.com/awslabs/eksdemo/releases/latest), look under Assets and locate the binary that matches your operation system and platform. Download the file, uncompress and copy to a location of your choice that is in your path. A common location on Mac and Linux is `/usr/local/bin`. Note that `eksctl` is required and [must be installed](https://docs.aws.amazon.com/eks/latest/userguide/eksctl.html) as well.

### Set the AWS Region

Most `eksdemo` commands require that you have configured a default AWS region or use the `--region` flag. There are 2 ways to configure a default region, either:

* Set in the the [AWS CLI Configuration file](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html). On Linux and MacOS this file is located in ~/.aws/config. You can use the `aws configure` command to set the region.
* Set the [`AWS_REGION` environment variable](https://docs.aws.amazon.com/sdkref/latest/guide/environment-variables.html) to the desired default region. An example is `export AWS_REGION=us-west-2`. Unless you set the environment variable in your `~/.bashrc` or `~/.zshrc`, you will need to set this every time you open a new terminal.

### Validate Install

To validate installation you can run the **`eksdemo version`** command and confirm you are running the latest version. The output will be similar to below:

```
» eksdemo version
eksdemo: version.Info{Version:"0.16.0", Date:"2024-08-19T17:41:55Z", Commit:"74fc767"}
```

To validate the AWS region is set, you can run **`eksdemo get cluster`** which will list running EKS clusters in the default region. If you don’t have any EKS clusters in the region, you will get the response: `No resources found.`

```
» eksdemo get cluster
+------------+--------+---------+---------+----------+----------+
|    Age     | Status | Cluster | Version | Platform | Endpoint |
+------------+--------+---------+---------+----------+----------+
| 3 weeks    | ACTIVE | green   |    1.27 | eks.17   | Public   |
| 20 minutes | ACTIVE | *blue   |    1.30 | eks.6    | Public   |
+------------+--------+---------+---------+----------+----------+
* Indicates current context in local kubeconfig
```

## Application Catalog

`eksdemo` comes with an extensive application catalog. Each application can be installed with a single command:
**`eksdemo install <application> -c <cluster-name> [flags]`**

To install applications under a group, you can use either a space or a hyphen. For example, each of the following are valid:
**`eksdemo install ingress nginx`** or **`eksdemo install ingress-nginx`**

The application catalog includes:

* `ack` — AWS Controllers for Kubernetes (ACK)
    * `apigatewayv2-controller` — ACK API Gateway v2 Controller
    * `ec2-controller` — ACK EC2 Controller
    * `ecr-controller` — ACK ECR Controller
    * `eks-controller` — ACK EKS Controller
    * `iam-controller` — ACK IAM Controller
    * `prometheusservice-controller` — ACK Prometheus Service Controller
    * `s3-controller` — ACK S3 Controller
* `adot-operator` — AWS Distro for OpenTelemetry (ADOT) Operator
* `appmesh-controller` — AWS App Mesh Controller
* `argo` — Get stuff done with Kubernetes!
    * `cd` — Declarative continuous deployment for Kubernetes
    * `workflows` — Workflow engine for Kubernetes
    * `workflows-cognito` — Workflow engine for Kubernetes using Cognito for authentication
* `autoscaling` — Kubernetes Autoscaling Applications
    * `cluster-autoscaler` — Kubernetes Cluster Autoscaler
    * `keda` — Kubernetes-based Event Driven Autoscaling
    * `vpa` — Vertical Pod Autoscaler
* `aws-fluent-bit` — AWS Fluent Bit
* `aws-lb-controller` — AWS Load Balancer Controller
* `cert-manager` — Cloud Native Certificate Management
* `cilium` — eBPF-based Networking, Observability, Security
* `consul` — HashiCorp Consul Service-Mesh
* `container-insights` — CloudWatch Container Insights
    * `adot-collector` — Container Insights ADOT Collector Metrics
    * `cloudwatch-agent` — Container Insights CloudWatch Agent Metrics
    * `fluent-bit` — Container Insights Fluent Bit Logs
    * `prometheus` — CloudWatch Container Insights monitoring for Prometheus
* `core-dump-handler` - Automatically saves core dumps to S3
* `crossplane` — The Cloud Native Control Plane
    * `core` — Crossplane Core Components
    * `ec2-provider` — Crossplane EC2 Provider
    * `iam-provider` — Crossplane IAM Provider
    * `s3-provider` — Crossplane S3 Provider
* `example` — Example Applications
    * `ascp` — Example for AWS Secrets Manager and Config Provider for Secret Store CSI Driver
    * `eks-workshop` — EKS Workshop Example Microservices
    * `game-2048` — Example Game 2048
    * `ghost` — Turn your audience into a business
    * `inflate` — Example App to Demonstrate Autoscaling
    * `kube-ops-view` — Kubernetes Operational View
    * `podinfo` — Go app w/microservices best practices
    * `wordpress` — WordPress Blog
* `external-dns` — ExternalDNS
* `falco` — Cloud Native Runtime Security
* `flux` — GitOps family of projects
    * `controllers` — Flux Controllers
    * `sync` — Flux GitRepository to sync with
* `goldilocks` — Get your resource requests "Just Right"
* `harbor` — Cloud Native Registry
* `headlamp` — An easy-to-use and extensible Kubernetes web UI
* `ingress` — Ingress Controllers
    * `contour` — Ingress Controller using Envoy proxy
    * `emissary` — Open Source API Gateway from Ambassador
    * `nginx` — Ingress NGINX Controller
* `istio` — Istio Service Mesh
    * `base` — Istio Base (includes CRDs)
    * `istiod` — Istio Control Plane
* `k8sgpt-operator` — K8sGPT Operator
* `karpenter` — Karpenter Node Autoscaling
* `keycloak-amg` — Keycloak SAML iDP for Amazon Managed Grafana
* `kube-prometheus` — End-to-end Cluster Monitoring with Prometheus
    * `karpenter-dashboards` — Karpenter Dashboards and ServiceMonitor
    * `stack` — Kube Prometheus Stack
    * `stack-amp` — Kube Prometheus Stack using Amazon Managed Prometheus
* `kube-state-metrics` — Kube State Metrics
* `kubecost` — Visibility Into Kubernetes Spend
    * `eks` — EKS optimized bundle of Kubecost
    * `eks-amp` — EKS optimized Kubecost using Amazon Managed Prometheus
    * `vendor` — Vendor distribution of Kubecost
* `metrics-server` — Kubernetes Metric Server
* `neuron` — AWS Neuron for Inferentia and Trainium Support
    * `neuron-device-plugin` — Neuron SDK Device Plugin
* `policy` — Kubernetes Policy Controllers
    * `kyverno` — Kubernetes Native Policy Management
    * `opa-gatekeeper` — Policy Controller for Kubernetes
* `prometheus-node-exporter` — Prometheus Node Exporter
* `secrets` — Secrets Management Solutions for Kubernetes
    * `store-csi-driver` — Integrates secrets stores with K8s via a CSI volume
    * `store-csi-driver-provider-aws` — AWS Secrets Manager and Config Provider for Secret Store CSI Driver
* `spark-operator` — Kubeflow Spark Operator
* `storage` — Kubernetes Storage Solutions
    * `ebs-csi` — Amazon EBS CSI driver
    * `efs-csi` — Amazon EFS CSI driver
    * `fsx-lustre-csi` — Amazon FSx for Lustre CSI Driver
    * `openebs` — Kubernetes storage simplified
* `vault` — HashiCorp Vault Secrets and Encryption Management System
* `velero` — Backup and Migrate Kubernetes Applications
* `vpc-lattice-controller` — Amazon VPC Lattice (Gateway API) Controller

## Kubectl-like get commands
`eksdemo` makes it easy to view AWS resources from the command line with commands that are very similar to how `kubectl get` works. Output defaults to a table, but raw AWS API output can be viewed with `-o yaml` and `-o json` flag options.

Almost all of the command have shorthand alaises to make it easier to type. For example, `get ec2` is an alias for `get ec2-instance`. You can find the aliases using the help command, `eksdemo get ec2-instance -h`.

* `access-entry` — EKS Access Entry
* `acm-certificate` — ACM Cerificate
* `addon` — EKS Managed Addon
* `addon-versions` — EKS Managed Addon Versions
* `alarm` — CloudWatch Alarm
* `amg-workspace` — Amazon Managed Grafana Workspace
* `ami` — Amazon Machine Image (AMI)
* `amp-rule` — Amazon Managed Prometheus Rule Namespace
* `amp-workspace` — Amazon Managed Prometheus Workspace
* `application` — Installed Applications
* `auto-scaling-group` — Auto Scaling Group
* `availability-zone` — Availability Zone
* `cloudformation-stack` — CloudFormation Stack
* `cloudtrail-event` — CloudTrail Event History
* `cloudtrail-trail` — CloudTrail Trail
* `cluster` — EKS Cluster
* `cognito` — Amazon Cognito Resources
    * `app-client` — Cognito User Pool App Client
    * `domain` — Cognito User Pool Domain
    * `user-pool` — Cognito User Pool
* `dns-record` — Route53 Resource Record Set
* `ec2-instance` — EC2 Instance
* `ecr-repository` — ECR Repository
* `elastic-ip` — Elastic IP Address
* `event-rule` — EventBridge Rule
* `fargate-profile` — EKS Fargate Profile
* `hosted-zone` — Route53 Hosted Zone
* `iam-oidc` — IAM OIDC Identity Provider
* `iam-policy` — IAM Policy
* `iam-role` — IAM Role
* `internet-gateway` — Internet Gateway
* `kms-key` — KMS Key
* `listener` — Load Balancer Listener
* `listener-rule` — Load Balancer Listener Rule
* `load-balancer` — Elastic Load Balancer
* `log-event` — CloudWatch Log Events
* `log-group` — CloudWatch Log Group
* `log-stream` — CloudWatch Log Stream
* `logs-insights` — CloudWatch Logs Insights
    * `query` —  Logs Insights Query History
    * `results` —  Logs Insights Query Results
    * `stats` —  Logs Insights Query Statistics
* `metric` — CloudWatch Metric
* `nat-gateway` — NAT Gateway
* `network-acl` — Network ACL
* `network-acl-rule` — Network ACL
* `network-interface` — Elastic Network Interface
* `node` — Kubernetes Node
* `nodegroup` — EKS Managed Nodegroup
* `organization` — AWS Organization
* `prefix-list` — Managed Prefix List
* `route-table` — Route Table
* `s3-bucket` — Amazon S3 Bucket
* `security-group` — Security Group
* `security-group-rule` — Security Group Rule
* `sqs-queue` — SQS Queue
* `ssm-node` — SSM Managed Node
* `ssm-parameter` — SSM Parameter
* `ssm-session` — SSM Session
* `subnet` — VPC Subnet
* `target-group` — Target Group
* `target-health` — Target Health
* `volume` — EBS Volume
* `vpc` — Virtual Private Cloud
* `vpc-endpoint` — VPC Endpoint
* `vpc-lattice` — VPC Lattice Resources
    * `service` —  VPC Lattice Service
    * `service-network` —  VPC Lattice Service Network
    * `target-group` —  VPC Lattice Target Group
* `vpc-summary` — VPC Summary

## Tutorials

The Basics tutorials provide detailed knowledge on how `eksdemo` works. It's recommended you review the Basics tutorials before proceeding to Advanced tutorial as they assume this knowlege.

### Basics
* [Create an Amazon EKS Cluster with Bottlerocket Nodes](/docs/create-cluster.md)
* [Request and Validate a Public Certificate with AWS Certificate Manager (ACM)](/docs/create-acm-cert.md)
* [Install AWS Load Balancer Controller](/docs/install-awslb.md)
* [Install ExternalDNS](/docs/install-edns.md)
* [Install Game 2048 Example Application with TLS using an ACM Certificate](/docs/install-game-2048.md)
* [Install Ingress NGINX](/docs/install-ingress-nginx.md)
* [Install cert-manager](/docs/install-cert-manager.md)
* [Install EBS CSI Driver](/docs/install-ebs-csi-driver.md)

### Advanced
* [Install Karpenter autoscaler and test node provisioning and consolidation](/docs/install-karpenter.md)
* [Install EKS optimized Kubecost using Amazon Managed Prometheus](/docs/install-kubecost.md)
* [Install Kube Prometheus Stack using Amazon Managed Prometheus](/docs/install-kube-prometheus.md)
* [Install Amazon VPC Lattice (Gateway API) Controller](/docs/install-vpc-lattice-controller.md)
* [Install HashiCorp Consul Application with Consul self-signed TLS](/docs/install-hashicorp-consul.md)

## Support & Feedback

This project is maintained by AWS Solution Architects. It is not part of an AWS service and support is provided best-effort by the maintainers. To post feedback, submit feature ideas, or report bugs, please use the [Issues](https://github.com/awslabs/eksdemo/issues) section of this repo. If you are interested in contributing, please see the [Contribution guide](CONTRIBUTING.md).

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.
