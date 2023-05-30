# Install ExternalDNS

ExternalDNS is a [Kubernetes SIGs](https://github.com/kubernetes-sigs) project that synchronizes exposed Kubernetes Services and Ingresses with DNS providers. It [watches the Kubernetes API](https://kubernetes.io/docs/reference/using-api/api-concepts/) for new `Service` and `Ingress` resources to determine which DNS records to configure.

1. [Prerequisites](#prerequisites)
2. [Install ExternalDNS](#install-externaldns-1)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install ExternalDNS

This section walks through the process of installing ExternalDNS. The command for performing the installation is:
**`eksdemo install external-dns -c <cluster-name>`**

Before you continue with the installation, you are encouraged to explore the help for external-dns with the `-h` flag. The exact syntax for the command is: **`eksdemo install external-dns -h`** 

Let's explore the dry run output with the `--dry-run` flag. The syntax for the command is: **`eksdemo install external-dns -c <cluster-name> --dry-run`**. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install external-dns -c <cluster-name> --dry-run
Creating 1 dependencies for external-dns
Creating dependency: external-dns-irsa

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
      name: external-dns
      namespace: external-dns
    roleName: eksdemo.blue.external-dns.external-dns
    roleOnly: true
    attachPolicy:
      Version: '2012-10-17'
      Statement:
      - Effect: Allow
        Action:
        - route53:ChangeResourceRecordSets
        Resource:
        - arn:aws:route53:::hostedzone/*
      - Effect: Allow
        Action:
        - route53:ListHostedZones
        - route53:ListResourceRecordSets
        Resource:
        - "*"


Helm Installer Dry Run:
+---------------------+------------------------------------------------+
| Application Version | v0.13.4                                        |
| Chart Version       | 1.12.2                                         |
| Chart Repository    | https://kubernetes-sigs.github.io/external-dns |
| Chart Name          | external-dns                                   |
| Release Name        | external-dns                                   |
| Namespace           | external-dns                                   |
| Wait                | false                                          |
+---------------------+------------------------------------------------+
Set Values: []
Values File:
---
image:
  tag: v0.13.4
provider: aws
registry: txt
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.external-dns.external-dns
  name: external-dns
txtOwnerId: blue
```

From the `--dry-run` output above, you can see that there is one dependency — an IAM Role. This role is associated with the ExternalDNS’s service account. It gives ExternalDNS permissions to query Hosted Zones and update Resource Record Sets. This is security best practices feature for EKS called [IAM Roles for Service Accounts (IRSA)](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html). `eksdemo` uses `eksctl` to create the IAM Role and the dry run output includes the exact configuration that will be used.

The values file configures ExternalDNS for the Helm chart install. It sets the provider to `aws`, uses TXT record entries to manage ownership of records and uses the cluster name as the record owner.


When you are ready to continue, proceed with installing ExternalDNS. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install external-dns -c <cluster-name>
Creating 1 dependencies for external-dns
Creating dependency: external-dns-irsa
2023-05-29 19:29:31 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/ebs-csi-controller-sa) will be excluded
2023-05-29 19:29:31 [ℹ]  1 iamserviceaccount (external-dns/external-dns) was excluded (based on the include/exclude rules)
2023-05-29 19:29:31 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-05-29 19:29:31 [ℹ]  no tasks
Downloading Chart: https://github.com/kubernetes-sigs/external-dns/releases/download/external-dns-helm-chart-1.12.2/external-dns-1.12.2.tgz
Helm installing...
2023/05/29 19:29:33 creating 1 resource(s)
2023/05/29 19:29:34 creating 5 resource(s)
Using chart version "1.12.2", installed "external-dns" version "v0.13.4" in namespace "external-dns"
NOTES:
***********************************************************************
* External DNS                                                        *
***********************************************************************
  Chart version: 1.12.2
  App version:   v0.13.4
  Image tag:     registry.k8s.io/external-dns/external-dns:v0.13.4
***********************************************************************
```

Let’s verify that the applications were installed properly with the **`eksdemo get application`** command. Since this command is specific to a given EKS cluster, the `-c <cluster-name>` flag is required. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get application -c <cluster-name>
+-------------------+--------------+---------+----------+--------+
|       Name        |  Namespace   | Version |  Status  | Chart  |
+-------------------+--------------+---------+----------+--------+
| aws-lb-controller | awslb        | v2.5.2  | deployed | 1.5.3  |
| external-dns      | external-dns | v0.13.4 | deployed | 1.12.2 |
+-------------------+--------------+---------+----------+--------+
```

From the output above we can see that both applications are successfully deployed in the EKS cluster named `blue`. `eksdemo` is using Helm as a Golang client library and the output above is very similar to running a `helm list --all-namespaces` command. Because Helm is bundled as a part of `eksdemo`, that you don’t need to have Helm installed on your system to install or manage any of the applications in the `eksdemo` application catalog.

When ExternalDNS is deployed on AWS, it will query Route 53 for a list of Hosted Zones. IAM Roles for Service Accounts (IRSA) is used to give permissions to access Route 53. You can quickly see all the IAM Roles configured for IRSA by using the **`eksdemo get iam-role -c <cluster-name>`** command. Include the `--last-used` or `-L` shorthand flag to see when the role was last used.

```
» eksdemo get iam-role -c <cluster-name> -L
+----------+-------------------------------------------------+------------+
|   Age    |                      Role                       | Last Used  |
+----------+-------------------------------------------------+------------+
| 14 hours | eksctl-blue-addon-vpc-cni-Role1-1PXCY1L5F2C05   | 1 hour     |
| 14 hours | eksdemo.blue.awslb.aws-load-balancer-controller | -          |
| 14 hours | eksdemo.blue.external-dns.external-dns          | 29 minutes |
| 14 hours | eksdemo.blue.karpenter.karpenter                | -          |
| 14 hours | eksdemo.blue.kube-system.cluster-autoscaler     | -          |
+----------+-------------------------------------------------+------------+
```

Notice that IAM Roles have been created for Cluster Autoscaler and Karpenter even though they haven’t been installed. See [Create an Amazon EKS Cluster with Bottlerocket Nodes](#create-an-amazon-eks-cluster-with-bottlerocket-nodes) for more detail on this and how to disable this feature.
