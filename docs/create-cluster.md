# Create an Amazon EKS Cluster with Bottlerocket Nodes

`eksdemo` can manage applications in any EKS cluster and the cluster doesn’t have to be created by `eksdemo`. You can use `eksctl` to create the cluster and then manage application using `eksdemo`. However, there are a number of benefits to using `eksdemo` to create your cluster:
* Cluster logging is enabled by default
* OIDC is enabled by default so IAM Roles for Service Accounts (IRSA) works out of the box
* The Managed Node Group ASG max is set to 10 so cluster autoscaling can work out of the box
* Private networking for nodes is set by default
* VPC CNI is configured as a Managed Add-on and configured with IRSA by default
* t3.large instances used by default instead of m5.large for cost savings, but can be easily changed with the `--instance` flag or the shorthand `-i`
* To use containerd as the CRI on Amazon EKS optimized Amazon Linux AMIs is as easy as using the `--containerd` boolean flag
* To create a Fargate profile that selects workloads in the “fargate” namespace, use the `--fargate` boolean flag
* Choose a supported EKS version with the `--version` flag or the shorthand `-v` like `-v 1.21`
* Using a different OS like Bottlerocket or Ubuntu is as easy as `--os bottlerocket` or `--os ubuntu`
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
      --containerd        use containerd runtime
      --dry-run           don't create, just print out all creation steps
      --fargate           create a Fargate profile
  -h, --help              help for cluster
  -i, --instance string   instance type (default "t3.large")
      --ipv6              use IPv6 networking
      --max int           max nodes (default 10)
      --min int           min nodes
      --no-roles          don't create IAM roles
  -N, --nodes int         desired number of nodes (default 2)
      --os string         Operating System (default "AmazonLinux2")
      --private           private cluster (includes ECR, S3, and other VPC endpoints)
  -v, --version string    Kubernetes version (default "1.24")

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
  version: "1.24"

addons:
- name: vpc-cni

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
      name: cluster-autoscaler
      namespace: kube-system
    roleName: eksdemo.blue.kube-system.cluster-autoscaler
    roleOnly: true
    attachPolicy:
      <snip>
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

managedNodeGroups:
- name: main
  amiFamily: Bottlerocket
  iam:
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy
    - arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
    - arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore
  instanceType: t3.xlarge
  minSize: 0
  desiredCapacity: 3
  maxSize: 10
  privateNetworking: true
  spot: false
```

You’ll notice that `eksdemo` automatically creates the IAM Roles used for IRSA for the most commonly deployed applications: AWS Load Balancer Controller, Cluster Autoscaler, ExternalDNS and Karpenter. This speeds up installation of the applications later as you don’t have to wait for CloudFormation to create the IAM Roles. To opt out, you can use the `--no-roles` flag.

After reviewing the output above, go ahead and create your cluster.

```
» eksdemo create cluster blue --os bottlerocket -i t3.xlarge -N 3
2023-01-25 09:04:24 [ℹ]  eksctl version 0.126.0
2023-01-25 09:04:24 [ℹ]  using region us-west-2
2023-01-25 09:04:24 [ℹ]  setting availability zones to [us-west-2d us-west-2b us-west-2a]
2023-01-25 09:04:24 [ℹ]  subnets for us-west-2d - public:192.168.0.0/19 private:192.168.96.0/19
2023-01-25 09:04:24 [ℹ]  subnets for us-west-2b - public:192.168.32.0/19 private:192.168.128.0/19
2023-01-25 09:04:24 [ℹ]  subnets for us-west-2a - public:192.168.64.0/19 private:192.168.160.0/19
2023-01-25 09:04:24 [ℹ]  nodegroup "main" will use "" [Bottlerocket/1.24]
2023-01-25 09:04:24 [ℹ]  using Kubernetes version 1.24
2023-01-25 09:04:24 [ℹ]  creating EKS cluster "blue" in "us-west-2" region with managed nodes
2023-01-25 09:04:24 [ℹ]  1 nodegroup (main) was included (based on the include/exclude rules)
<snip>
2023-01-25 09:23:26 [ℹ]  waiting for CloudFormation stack "eksctl-blue-nodegroup-main"
2023-01-25 09:23:26 [ℹ]  waiting for the control plane to become ready
2023-01-25 09:23:28 [✔]  saved kubeconfig as "/Users/awsuser/.kube/config"
2023-01-25 09:23:28 [ℹ]  no tasks
2023-01-25 09:23:28 [✔]  all EKS cluster resources for "blue" have been created
2023-01-25 09:23:29 [ℹ]  kubectl command should work with "/Users/awsuser/.kube/config", try 'kubectl get nodes'
2023-01-25 09:23:29 [✔]  EKS cluster "blue" in "us-west-2" region is ready
```

To view the status and info about your cluster you can run the **`eksdemo get cluster`** command.

```
» eksdemo get cluster
+------------+--------+---------+---------+----------+----------+---------+
|    Age     | Status | Cluster | Version | Platform | Endpoint | Logging |
+------------+--------+---------+---------+----------+----------+---------+
| 3 weeks    | ACTIVE | green   |    1.23 | eks.5    | Public   | true    |
| 20 minutes | ACTIVE | *blue   |    1.24 | eks.3    | Public   | true    |
+------------+--------+---------+---------+----------+----------+---------+
* Indicates current context in local kubeconfig
```

To view detail on the node group, use the **`eksdemo get nodegroup`** command. For this get command and many others, there is a required `--cluster <cluster-name>` flag.

```
» eksdemo get nodegroup --cluster blue
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
|    Age    | Status | Name | Nodes | Min | Max |     Version     |   Type    | Instance(s) |
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
| 5 minutes | ACTIVE | main |     3 |   0 |  10 | 1.11.1-104f8e0f | ON_DEMAND | t3.xlarge   |
+-----------+--------+------+-------+-----+-----+-----------------+-----------+-------------+
```

To view detail on the nodes, use the **`eksdemo get node`** command. Here we’ll use the `-c` flag which is the shorthand version of the `--cluster` flag.

```
» eksdemo get node -c blue
+-----------+----------------------+---------------------+-----------+------------+-----------+
|    Age    |         Name         |     Instance Id     |   Type    |    Zone    | Nodegroup |
+-----------+----------------------+---------------------+-----------+------------+-----------+
| 5 minutes | ip-192-168-112-160.* | i-01049dccf2e58d265 | t3.xlarge | us-west-2d | main      |
| 5 minutes | ip-192-168-141-119.* | i-003139b73a29ff1b7 | t3.xlarge | us-west-2b | main      |
| 5 minutes | ip-192-168-186-156.* | i-0583cab4366088ac2 | t3.xlarge | us-west-2a | main      |
+-----------+----------------------+---------------------+-----------+------------+-----------+
* Names end with "us-west-2.compute.internal"
```

Congratulations, your EKS cluster with 3 Bottlerocket `t3.xlarge` nodes is now ready!  In the future if you want to see more detail from a get command you can use `-o yaml` or `-o json` and you will see the raw AWS API response in full. For example, you can try running **`eksdemo get cluster blue -o yaml`**. You can also run **`eksdemo get`** to see all the options available.
