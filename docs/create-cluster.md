# Create an Amazon EKS Cluster with Bottlerocket Nodes

`eksdemo` can manage applications in any EKS cluster and the cluster doesn’t have to be created by `eksdemo`. You can use `eksctl` to create the cluster and then manage application using `eksdemo`. However, there are a number of benefits to using `eksdemo` to create your cluster:
* Cluster logging is enabled by default
* OIDC is enabled by default so IAM Roles for Service Accounts (IRSA) works out of the box
* The Managed Node Group ASG max is set to 10 so Cluster Autoscaler can work out of the box
* Private networking for nodes is set by default
* VPC CNI is configured as a Managed Add-on and configured with IRSA, with network policy enabled
* t3.large instances used by default instead of m5.large for cost savings, but can be easily changed with the `--instance` flag or the shorthand `-i`
* To create a Fargate profile that selects workloads in the “fargate” namespace, use the `--fargate` boolean flag
* Choose a supported EKS version with the `--version` flag or the shorthand `-v` like `-v 1.29`
* Using a different OS like Bottlerocket or AL2023 is as easy as `--os bottlerocket` or `--os amazonlinux2023`
* To use IPv6 networking, set the `--ipv6` boolean flag
* If you need to further customize the config, add the `--dry-run` flag and it will output the eksctl YAML config file and you can copy/paste it into a file, make your edits and run `eksctl create cluster -f cluster.yaml`

In this section we will walk through the process of creating an Amazon EKS cluster using `eksdemo` that highlights some of the benefits from the list above. First, review the usage and options of the **`eksdemo create cluster`** command using the help flag `--help` or the shorthand `-h`.

```
» eksdemo create cluster -h
Create EKS Cluster

Usage:
  eksdemo create cluster NAME [flags]

Aliases:
  cluster, clusters

Flags:
      --disable-network-policy   don't enable network policy for Amazon VPC CNI
      --dry-run                  don't create, just print out all creation steps
      --encrypt-secrets string   alias of KMS key to encrypt secrets
      --fargate                  create a Fargate profile
  -h, --help                     help for cluster
  -H, --hostname-type string     type of hostname to use for EC2 instances (default "resource-name")
  -i, --instance string          instance type (default "t3.large")
      --ipv6                     use IPv6 networking
      --kubeconfig string        path to write kubeconfig (default "/Users/jsmith/.kube/config")
      --max int                  max nodes (default 10)
      --min int                  min nodes
      --no-roles                 don't create IAM roles
      --no-taints                don't taint nodes with GPUs or Neuron cores
  -N, --nodes int                desired number of nodes (default 2)
      --os string                Operating System (default "AmazonLinux2")
      --prefix-assignment        configure VPC CNI for prefix assignment
      --private                  private cluster (includes ECR, S3, and other VPC endpoints)
  -v, --version string           Kubernetes version (default "1.30")
      --vpc-cidr string          CIDR to use for EKS Cluster VPC (default "192.168.0.0/16")
      --zones strings            list of AZs to use. ie. us-east-1a,us-east-1b,us-east-1c

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
  ```

In this example, we would like the following customizations:
* Name our cluster “blue”
* Use Bottlerocket nodes instead of Amazon Linux 2
* Use `t3.xlarge` instances instead t3.large
* Create a Managed Node Group with 3 nodes instances instead of 2

The command for this is **`eksdemo create cluster blue --os bottlerocket -i t3.xlarge -N 3`**

Before you run the command, let’s dive a bit deeper and understand exactly how `eksdemo` will use and configure `eksctl` to create the cluster. We can do that with the `--dry-run` flag.

```
» eksdemo create cluster blue --os bottlerocket -i t3.xlarge -N 3 --dry-run

Eksctl Resource Manager Dry Run:
eksctl create cluster -f -
---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: blue
  region: us-west-2
  version: "1.30"

addons:
- name: vpc-cni
  version: latest
  configurationValues: |-
    enableNetworkPolicy: "true"
    env:
      ENABLE_PREFIX_DELEGATION: "false"

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]

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
  - metadata:
      name: ebs-csi-controller-sa
      namespace: kube-system
    roleName: eksdemo.blue.kube-system.ebs-csi-controller-sa
    roleOnly: true
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy
  - metadata:
      name: external-dns
      namespace: external-dns
    roleName: eksdemo.blue.external-dns.external-dns
    roleOnly: true
    attachPolicy:
      <snip>
  - metadata:
      name: karpenter
      namespace: karpenter
    roleName: eksdemo.blue.karpenter.karpenter
    roleOnly: true
    attachPolicy:
      <snip>

vpc:
  cidr: 192.168.0.0/16
  hostnameType: resource-name

managedNodeGroups:
- name: main
  amiFamily: Bottlerocket
  desiredCapacity: 3
  iam:
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
    - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
    - arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore
  instanceType: t3.xlarge
  minSize: 0
  maxSize: 10
  privateNetworking: true
  spot: false
```

You’ll notice that `eksdemo` automatically creates the IAM Roles used for IRSA for the most commonly deployed applications: AWS Load Balancer Controller, Cluster Autoscaler, ExternalDNS and Karpenter. This speeds up installation of the applications later as you don’t have to wait for CloudFormation to create the IAM Roles. To opt out, you can use the `--no-roles` flag.

After reviewing the output above, go ahead and create your cluster.

```
» eksdemo create cluster blue --os bottlerocket -i t3.xlarge -N 3
2024-07-30 18:52:39 [ℹ]  eksctl version 0.180.0
2024-07-30 18:52:39 [ℹ]  using region us-west-2
2024-07-30 18:52:39 [ℹ]  setting availability zones to [us-west-2c us-west-2a us-west-2d]
2024-07-30 18:52:39 [ℹ]  subnets for us-west-2c - public:192.168.0.0/19 private:192.168.96.0/19
2024-07-30 18:52:39 [ℹ]  subnets for us-west-2a - public:192.168.32.0/19 private:192.168.128.0/19
2024-07-30 18:52:39 [ℹ]  subnets for us-west-2d - public:192.168.64.0/19 private:192.168.160.0/19
2024-07-30 18:52:39 [ℹ]  nodegroup "main" will use "" [Bottlerocket/1.30]
2024-07-30 18:52:39 [ℹ]  using Kubernetes version 1.30
2024-07-30 18:52:39 [ℹ]  creating EKS cluster "blue" in "us-west-2" region with managed nodes
2024-07-30 18:52:39 [ℹ]  1 nodegroup (main) was included (based on the include/exclude rules)
<snip>
2024-07-30 19:09:37 [ℹ]  waiting for CloudFormation stack "eksctl-blue-nodegroup-main"
2024-07-30 19:09:37 [ℹ]  waiting for the control plane to become ready
2024-07-30 19:09:38 [✔]  saved kubeconfig as "/Users/awsuser/.kube/config"
2024-07-30 19:09:38 [ℹ]  no tasks
2024-07-30 19:09:38 [✔]  all EKS cluster resources for "blue" have been created
2024-07-30 19:09:38 [✔]  created 1 managed nodegroup(s) in cluster "blue"
2024-07-30 19:09:40 [ℹ]  kubectl command should work with "/Users/awsuser/.kube/config", try 'kubectl get nodes'
2024-07-30 19:09:40 [✔]  EKS cluster "awsuser" in "us-west-2" region is ready
```

To view the status and info about your cluster you can run the **`eksdemo get cluster`** command.

```
» eksdemo get cluster
+------------+--------+---------+---------+----------+----------+
|    Age     | Status | Cluster | Version | Platform | Endpoint |
+------------+--------+---------+---------+----------+----------+
| 3 weeks    | ACTIVE | green   |    1.28 | eks.16   | Public   |
| 20 minutes | ACTIVE | *blue   |    1.30 | eks.5    | Public   |
+------------+--------+---------+---------+----------+----------+
* Indicates current context in local kubeconfig
```

To view detail on the node group, use the **`eksdemo get nodegroup`** command. For this get command and many others, there is a required `--cluster <cluster-name>` flag.

```
» eksdemo get nodegroup --cluster blue
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
|    Age    | Status | Name | Nodes | Min | Max |     Version     |   Type    | Instance(s) |
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
| 5 minutes | ACTIVE | main |     3 |   0 |  10 | 1.20.5-a3e8bda1 | ON_DEMAND | t3.xlarge   |
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
```

To view detail on the nodes, use the **`eksdemo get node`** command. Here we’ll use the `-c` flag which is the shorthand version of the `--cluster` flag.

```
» eksdemo get node -c blue
+-----------+-----------------------+---------------------+-----------+------------+-----------+
|    Age    |         Name          |     Instance Id     |   Type    |    Zone    | Nodegroup |
+-----------+-----------------------+---------------------+-----------+------------+-----------+
| 5 minutes | i-058de3c37e4d56968.* | i-058de3c37e4d56968 | t3.xlarge | us-west-2b | main      |
| 5 minutes | i-05e74a812e705a2b4.* | i-05e74a812e705a2b4 | t3.xlarge | us-west-2c | main      |
| 5 minutes | i-0d14753576296c6e0.* | i-0d14753576296c6e0 | t3.xlarge | us-west-2a | main      |
+-----------+-----------------------+---------------------+-----------+------------+-----------+
* Names end with "us-west-2.compute.internal"
```

Congratulations, your EKS cluster with 3 Bottlerocket `t3.xlarge` nodes is now ready!  In the future if you want to see more detail from a get command you can use `-o yaml` or `-o json` and you will see the raw AWS API response in full. For example, you can try running **`eksdemo get cluster blue -o yaml`**. You can also run **`eksdemo get`** to see all the options available.
