# Install AWS Load Balancer Controller

The AWS Load Balancer Controller manages Elastic Load Balancers for Kubernetes clusters. The controller provision a Network Load Balancer (NLB) when you create a Kubernetes service of type `LoadBalancer` and provisions an Application Load Balancer (ALB) when you create a Kubernetes `Ingress`. 

1. [Prerequisites](#prerequisites)
2. [Install the AWS Load Balancer Controller](#install-the-aws-load-balancer-controller)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install the AWS Load Balancer Controller

In this section we walk through the process of installing the AWS Load Balancer Controller application. The command for performing the installation is: **`eksdemo install aws-lb-controller -c <cluster-name>`**

Let’s learn a bit more about the command and it’s options before we continue by using the `-h` help shorthand flag.

```
» eksdemo install aws-lb-controller -h
Install aws-lb-controller

Usage:
  eksdemo install aws-lb-controller [flags]

Aliases:
  aws-lb-controller, aws-lb, awslb

Flags:
      --chart-version string     chart version (default "1.4.7")
  -c, --cluster string           cluster to install application (required)
      --default                  set as the default IngressClass for the cluster
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for aws-lb-controller
  -n, --namespace string         namespace to install (default "awslb")
      --service-account string   service account name (default "aws-load-balancer-controller")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous             use previous working chart/app versions ("1.4.6"/"v2.4.5")
  -v, --version string           application version (default "v2.4.6")
```

The help content provides a lot of valuable information at a glance. The default chart and application versions, namespace and service account names are included along with optional flags to modify the defaults if desired.

A very powerful optional flag is the `--dry-run` flag. This will print out details about any dependencies and exactly how the application install will take place so there is no mystery about the steps `eksdemo` is taking to install your application. Let’s use the the `--dry-run` flag to understand how the AWS Load Balancer Controller will be installed. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install aws-lb-controller -c <cluster-name> --dry-run
Creating 1 dependencies for aws-lb-controller
Creating dependency: aws-lb-controller-irsa

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
      name: aws-load-balancer-controller
      namespace: awslb
    roleName: eksdemo.blue.awslb.aws-load-balancer-controller
    roleOnly: true
    attachPolicy:
      <snip>

Helm Installer Dry Run:
+---------------------+----------------------------------+
| Application Version | v2.4.6                           |
| Chart Version       | 1.4.7                            |
| Chart Repository    | https://aws.github.io/eks-charts |
| Chart Name          | aws-load-balancer-controller     |
| Release Name        | aws-lb-controller                |
| Namespace           | awslb                            |
| Wait                | false                            |
+---------------------+----------------------------------+
Set Values: []
Values File:
---
replicaCount: 1
image:
  tag: v2.4.6
fullnameOverride: aws-load-balancer-controller
clusterName: blue
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.awslb.aws-load-balancer-controller
  name: aws-load-balancer-controller
region: us-west-2
vpcId: vpc-08a68dc8b440fec75
```

From the `--dry-run` output above, you can see that there is one dependency — an IAM Role. This role is associated with the AWS Load Balancer Controller’s service account. This is security best practices feature for EKS called [IAM Roles for Service Accounts (IRSA)](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html). `eksdemo` uses `eksctl` to create the IAM Role and the dry run output includes the exact configuration that will be used.

Additionally, the output includes details on how the application will be installed. Most applications, including the AWS Load Balancer Controller, are installed using a Helm chart. The dry run information for Helm installs includes 3 sections:

1. A table with the Helm chart repository URL and name, the chart and application versions and the release name and namespace where the application will be installed.
2. Any `--set` flag variables to override the chart’s `values.yaml` defaults or the values file configuration in the next section. See the Helm documentation for more details on the [format and limitations of the --set flag](https://helm.sh/docs/intro/using_helm/#the-format-and-limitations-of---set).
3. The opinionated values file settings built into the `eksdemo` application catalog. Some of these settings can be change with optional flags. If a flag is not available for the value you wish to change, the `--set` flag can be used to override any value.

With this application and with many others, a number of values file settings are automatically populated. In the example above, the `region`, `vpcID` and AWS Account ID in the IRSA annotation are dynamically updated each time `eksdemo` runs.

With a thorough understanding of how the application install process works, let’s install the AWS Load Balancer controller. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install aws-lb-controller -c <cluster-name>
Creating 1 dependencies for aws-lb-controller
Creating dependency: aws-lb-controller-irsa
2023-01-30 21:49:52 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/cluster-autoscaler) will be excluded
2023-01-30 21:49:52 [ℹ]  1 iamserviceaccount (awslb/aws-load-balancer-controller) was excluded (based on the include/exclude rules)
2023-01-30 21:49:52 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-01-30 21:49:52 [ℹ]  no tasks
Downloading Chart: https://aws.github.io/eks-charts/aws-load-balancer-controller-1.4.7.tgz
Helm installing...
2023/01/30 21:49:54 creating 2 resource(s)
2023/01/30 21:49:54 Clearing discovery cache
2023/01/30 21:49:54 beginning wait for 2 resources with timeout of 1m0s
2023/01/30 21:50:02 creating 1 resource(s)
2023/01/30 21:50:03 creating 12 resource(s)
Using chart version "1.4.7", installed "aws-lb-controller" version "v2.4.6" in namespace "awslb"
NOTES:
AWS Load Balancer controller installed!
```