# Install and Test Karpenter Autoscaling

With `eksdemo`, it’s very easy to setup, install and test Karpenter autoscaling. Karpenter can be installed with a single command and will automate the following pre-install and post-install steps: 

1. Create the EC2 Spot Service Linked Role (if it doesn't already exist)
2. Create the Karpenter Controller IAM Role (IRSA)
3. Create the Karpenter Node IAM Role
4. Create an SQS Queue and EventBridge rules for native Spot Termination Handling
5. Add an entry to the `aws-auth` ConfigMap for the Karpenter Node IAM Role
6. Install the Karpenter Helm Chart
7. Create default Karpenter `NodePool` and `EC2NodeClass` Custom Resources

This tutorial walks through the installation of the Karpenter Autoscaler and the example Inflate application to trigger an autoscaling event. It also tests Node Consolidation.

1. [Prerequisites](#prerequisites)
2. [Install Karpenter Autoscaler](#install-karpenter-autoscaler)
3. [Test Automatic Node Provisioning](#test-automatic-node-provisioning)
4. [Test Node Consolidation](#test-node-consolidation)
5. [(Optional) Install Karpenter Dashboards](#optional-install-karpenter-dashboards)
6. [(Optional) Inspect Karpenter IAM Roles](#optional-inspect-karpenter-iam-roles)
7. [(Optional) Inspect Karpenter SQS Queue and EventBridge Rules](#optional-inspect-karpenter-sqs-queue-and-eventbridge-rules)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.
```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install Karpenter Autoscaler

In this section we walk through the process of installing the Karpenter Autoscaler. The command for performing the installation is: `eksdemo install karpenter -c <cluster-name>`

Let’s explore the command and its options by using the -h help shorthand flag.

```
» eksdemo install karpenter -h
Install karpenter

Usage:
  eksdemo install karpenter [flags]

Flags:
  -a, --ami-family string          node class AMI family (default "AL2")
      --chart-version string       chart version (default "1.0.0")
  -c, --cluster string             cluster to install application (required)
      --consolidate-after string   time after a pod is scheduled/removed before considering the node consolidatable (default "1m")
      --dry-run                    don't install, just print out all installation steps
      --enable-spottospot          enables the spot to spot consolidation feature
      --expire-after string        time a node can live on the cluster before being deleted (default "720h")
  -h, --help                       help for karpenter
  -n, --namespace string           namespace to install (default "karpenter")
      --replicas int               number of replicas for the controller deployment (default 1)
      --service-account string     service account name (default "karpenter")
      --set strings                set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous               use previous working chart/app versions ("1.0.0"/"1.0.0")
  -v, --version string             application version (default "1.0.0")
```

The Karpenter specific flags are:
* `--ami-family` -- This sets the AMI family on the default `EC2NodeClass`. Options include AL2, AL2023, and Bottlerocket.
* `--consolidate-after` -- This indicate how long Karpenter should wait after a pod schedules or is removed from the node before considering the node consolidatable.
* `--enable-spottospot` -- As of v0.34.x, [SpotToSpotConsolidation](https://karpenter.sh/docs/concepts/disruption/#spot-consolidation) is in ALPHA and disabled by default. This flag enables this feature.
* `--expire-after` -- Karpenter will mark nodes as expired and disrupt them after they have lived a set number of seconds. The default of 720 hours is 30 days.
* `--replicas` -- `eksdemo` defaults to only 1 replica for easier log viewing in a demo environment. You can use this flag to increase to the default Karpenter Helm chart value of 2 replicas for high availability.

The `eksdemo` install of Karpenter is identical to the [Getting Started with eksctl](https://karpenter.sh/docs/getting-started/getting-started-with-eksctl/) instructions with the following small differences:
* The Karpenter instructions create a cluster with `karpenter.sh/discovery: <cluster-name>` tags. This is not required and any `eksctl` or `eksdemo` created EKS cluster will work.
* The IRSA role, Node IAM Role and SQS Queue with EventBridge rules are in separate CloudFormation stacks. The IRSA role is created using `eksctl`.
* The SQS Queue name is `karpenter-<cluster-name>` instead of `<cluster-name>`.
* The four EventBridge rule names follow the pattern `karpenter-<cluster-name>-<rule-name>` instead of the automated CloudFormation generated names.
* The Karpenter controller deployment defaults to 1 replica instead of 2.
* There are a few changes to the  default Karpenter `NodePool` and `EC2NodeClass` Custom Resources:
    * Support for both `on-demand` and `spot` nodes is configured instead of only `spot`. It will still default to using Spot nodes but you can request `on-demand` using the `karpenter.sh/capacity-type` label selector.
    * The `subnetSelector` uses `eksctl` default tags to select only Private subnets.
    * The `securityGroupSelector` uses the `aws:eks:cluster-name: <cluster-name>` tag selector to use the EKS cluster security group.

Optionally, if you want to see details about all the prerequisite items that are created when you install Karpenter, you can run the Karpenter install command with the `--dry-run` flag to first inspect all the actions that will be performed. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install karpenter -c <cluster-name> --dry-run
Creating 5 dependencies for karpenter
<snip>
Helm Installer Dry Run:
+---------------------+------------------------------------------+
| Application Version | 1.0.0                                    |
| Chart Version       | 1.0.0                                    |
| Chart Repository    | oci://public.ecr.aws/karpenter/karpenter |
| Chart Name          | karpenter                                |
| Release Name        | karpenter                                |
| Namespace           | karpenter                                |
| Wait                | true                                     |
+---------------------+------------------------------------------+
Set Values: []
Values File:
---
serviceAccount:
  name: karpenter
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.karpenter.karpenter
replicas: 1
controller:
  image:
    tag: 1.0.0
  resources:
    requests:
      cpu: "1"
      memory: "1Gi"
settings:
  clusterName: blue
  interruptionQueue: karpenter-blue
  featureGates:
    # -- spotToSpotConsolidation is ALPHA and is disabled by default.
    # Setting this to true will enable spot replacement consolidation for both single and multi-node consolidation.
    spotToSpotConsolidation: false

Creating 1 post-install resources for karpenter
Creating post-install resource: karpenter-default-nodepool

Kubernetes Resource Manager Dry Run:
---
apiVersion: karpenter.sh/v1
kind: NodePool
metadata:
  name: default
spec:
  template:
    spec:
      requirements:
        - key: kubernetes.io/arch
          operator: In
          values: ["amd64"]
        - key: kubernetes.io/os
          operator: In
          values: ["linux"]
        - key: karpenter.sh/capacity-type
          operator: In
          values: ["on-demand", "spot"]
        - key: karpenter.k8s.aws/instance-category
          operator: In
          values: ["c", "m", "r"]
        - key: karpenter.k8s.aws/instance-generation
          operator: Gt
          values: ["2"]
      nodeClassRef:
        group: karpenter.k8s.aws
        kind: EC2NodeClass
        name: default
      expireAfter: 720h
  limits:
    cpu: 1000
  disruption:
    consolidationPolicy: WhenEmptyOrUnderutilized
    consolidateAfter: 1m
---
apiVersion: karpenter.k8s.aws/v1
kind: EC2NodeClass
metadata:
  name: default
spec:
  amiFamily: AL2
  role: KarpenterNodeRole-blue
  subnetSelectorTerms:
    - tags:
        Name: eksctl-blue-cluster/SubnetPrivate*
  securityGroupSelectorTerms:
    - tags:
        aws:eks:cluster-name: blue
  tags:
    eksdemo.io/version: v0.16.0
  amiSelectorTerms:
    - id: ami-0d6c6ac554d62a745
    - id: ami-0eb2ae1eda347f8cf
```

Now, install Karpenter.
```
» eksdemo install karpenter -c <cluster-name>
Creating 5 dependencies for karpenter
<snip>
Downloading Chart: oci://public.ecr.aws/karpenter/karpenter:1.0.0
Helm installing...
2024/08/17 17:45:41 creating 1 resource(s)
2024/08/17 17:45:42 creating 1 resource(s)
2024/08/17 17:45:42 creating 1 resource(s)
2024/08/17 17:45:42 beginning wait for 3 resources with timeout of 1m0s
2024/08/17 17:45:43 Clearing REST mapper cache
2024/08/17 17:45:45 creating 1 resource(s)
2024/08/17 17:45:46 creating 16 resource(s)
2024/08/17 17:45:47 beginning wait for 16 resources with timeout of 5m0s
<snip>
2024/08/17 17:46:02 beginning wait for 1 resources to be deleted with timeout of 5m0s
Using chart version "1.0.0", installed "karpenter" version "1.0.0" in namespace "karpenter"
Creating 1 post-install resources for karpenter
Creating post-install resource: karpenter-default-nodepool
Creating NodePool "default"
Creating EC2NodeClass "default"
```

## Test Automatic Node Provisioning

The test of automatic node provisioning is identical to the [6. Scale up deployment](https://karpenter.sh/docs/getting-started/getting-started-with-karpenter/#6-scale-up-deployment) Getting Started instructions. `eksdemo` uses the same Deployment manifest but uses Helm to install it so it can be managed and uninstalled. Additionally there are a few options. Let's review the help.

```
» eksdemo install example-inflate -h
Install example-inflate

Usage:
  eksdemo install example-inflate [flags]

Flags:
  -c, --cluster string     cluster to install application (required)
      --dry-run            don't install, just print out all installation steps
  -h, --help               help for example-inflate
  -n, --namespace string   namespace to install (default "inflate")
      --on-demand          request on-demand instances using karpenter node selector
      --replicas int       number of replicas for the deployment (default 0)
      --spread             use topology spread constraints to spread across zones
```
The Inflate specific flags are:
* `--on-demand` -- Requests On-Demand nodes by adding a Node Selector for the `karpenter.sh/capacity-type: on-demand` label 
* `--replicas` -- Sets the number of replicas at install time instead of scaling after install
* `--spread` -- Spreads the replicas across availability zones using Topology Spread Constraint and `topologyKey: topology.kubernetes.io/zone`

Let's test autoscaling with 10 replicas and spread across AZs. We will first review the manifest with the `--dry-run` flag. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install example-inflate -c <cluster-name> --replicas 10 --spread --dry-run

Manifest Installer Dry Run:
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inflate
  namespace: inflate
spec:
  replicas: 10
  selector:
    matchLabels:
      app: inflate
  template:
    metadata:
      labels:
        app: inflate
    spec:
      terminationGracePeriodSeconds: 0
      securityContext:
        runAsUser: 1000
        runAsGroup: 3000
        fsGroup: 2000
      containers:
      - name: inflate
        image: public.ecr.aws/eks-distro/kubernetes/pause:3.7
        resources:
          requests:
            cpu: 1
        securityContext:
          allowPrivilegeEscalation: false
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: inflate
```

Now, remove the `--dry-run` flag and install the Inflate app to trigger an autoscaling event.

```
» eksdemo install example-inflate -c <cluster-name> --replicas 10 --spread
Helm installing...
2024/06/06 16:31:26 creating 1 resource(s)
2024/06/06 16:31:26 creating 1 resource(s)
Using chart version "n/a", installed "example-inflate" version "n/a" in namespace "inflate"
```

[Optional] If you want an interactive way of watching the workload changes, install the EKS Node Viewer from here: https://github.com/awslabs/eks-node-viewer

After that it should show the output of eks-node-viewer like below:
![eks node viewer workload changes](/docs/images/eks_node_viewer1.png?raw=true "eks_node_viewer_workload_changes")

Wait a few moments and then list the EC2 instances in your EKS cluster's VPC using the `eksdemo get ec2-instance` command. The `-c` shorthand cluster flag is optional and filters the instance list to your EKS Cluster VPC.

```
» eksdemo get ec2-instance -c <cluster-name>
+-----------+---------+---------------------+--------------------------------+--------------+------------+
|    Age    |  State  |         Id          |              Name              |     Type     |    Zone    |
+-----------+---------+---------------------+--------------------------------+--------------+------------+
| 5 hours   | running | i-0700781fbdfc33254 | blue-main-Node                 | t3.large     | us-west-2d |
| 1 minute  | running | i-0b9d97609d89c4e87 | karpenter.sh/provisioner-na... | *c5ad.xlarge | us-west-2d |
| 5 hours   | running | i-0f4aef0997638970d | blue-main-Node                 | t3.large     | us-west-2b |
| 1 minute  | running | i-01833ebd5a38c541c | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2b |
| 1 minute  | running | i-0fdc7907aab2654ea | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2c |
+-----------+---------+---------------------+--------------------------------+--------------+------------+
* Indicates Spot Instance
```

In the example above, you can see that Karpenter has created 3 Spot instances all in separate availability zones: `us-west-2b`, `us-west-2c` and `us-west-2d`. The 10 pods are spread across all 3 AZ's due to the topology spread constraints.

## Test Node Consolidation

To test consolidation, let's reduce the number of replicas for the Inflate deployment from 10 to 5.

```
» kubectl -n inflate scale deploy/inflate --replicas 5
deployment.apps/inflate scaled
```

You may need to wait a few minutes for Karpenter's consolidation logic to make a decision to replace or terminate nodes. Again, use the `eksdemo get ec2-instance` command to view the EC2 instances. This time we will use the shorthand alias `ec2`.


If you have installed the eks-node-viewer, the output would be like below:
![eks_node_viewer_workload_changes](/docs/images/eks_node_viewer2.png?raw=true "eks_node_viewer_workload_changes")
```
» eksdemo get ec2 -c <cluster-name>
+------------+------------+---------------------+--------------------------------+--------------+------------+
|    Age     |   State    |         Id          |              Name              |     Type     |    Zone    |
+------------+------------+---------------------+--------------------------------+--------------+------------+
| 6 hours    | running    | i-0700781fbdfc33254 | blue-main-Node                 | t3.large     | us-west-2d |
| 11 minutes | running    | i-0b9d97609d89c4e87 | karpenter.sh/provisioner-na... | *c5ad.xlarge | us-west-2d |
| 6 hours    | running    | i-0f4aef0997638970d | blue-main-Node                 | t3.large     | us-west-2b |
| 11 minutes | terminated | i-01833ebd5a38c541c | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2b |
| 11 minutes | running    | i-0fdc7907aab2654ea | karpenter.sh/provisioner-na... | *m3.xlarge   | us-west-2c |
+------------+------------+---------------------+--------------------------------+--------------+------------+
* Indicates Spot Instance
```

In the example above, Karpenter decided to terminate one of the Spot instances in `us-west-2b` because there is already a Managed Node Group node running in `us-west-2b` so the deployment is still spread across AZ's.

## (Optional) Install Karpenter Dashboards

The Karpenter team has created 2 Grafana dashboards that provide a variety of visualization examples of Karpenter metrics:
* Karpenter Capacity Dashboard
* Karpenter Performance Dashboard

If you already have Prometheus Operator installed, you can [Install the Karpenter Dashboards](/docs/install-kube-prometheus.md#optional-install-karpenter-dashboards). Otherwise you can use `eksdemo` to [Install the Kube Prometheus Stack](/docs/install-kube-prometheus.md) first, which includes Prometheus Operator.

![Karpenter Performance Dashboard](/docs/images/karpenter-perf-dashboard.png "Karpenter Performance Dashboard")

## (Optional) Inspect Karpenter IAM Roles 

The Karpenter install process creates 2 IAM Roles:
* Karpenter Controller IRSA (IAM Role for Service Accounts)
* Karpenter Node IAM Role

Use the `eksdemo get iam-role` command and the `--search` flag to find the roles.

```
» eksdemo get iam-role --search karpenter
+---------+----------------------------------+
|   Age   |              Role                |
+---------+----------------------------------+
| 7 hours | eksdemo.blue.karpenter.karpenter |
| 6 hours | KarpenterNodeRole-blue           |
+---------+----------------------------------+
```

`eksdemo` uses a specific naming convention for IRSA roles: `eksdemo.<cluster-name>.<namespace>.<serviceaccount-name>`. To view the permissions assigned to the role, use the `eksdemo get iam-policy` command and the `--role` command which lists only the policies assigned to the role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role eksdemo.<cluster-name>.karpenter.karpenter
+-----------------------------------------------------------------+--------+-------------+
|                               Name                              |  Type  | Description |
+-----------------------------------------------------------------+--------+-------------+
| eksctl-blue-addon-iamserviceaccount-karpenter-karpenter-Policy1 | Inline |             |
+-----------------------------------------------------------------+--------+-------------+
```

To view the details of the policy, including the policy document, you can use the `--output` flag or `-o` shorthand flag to output the raw AWS API responses in either JSON or YAML. In the example we'll use YAML.

```
» eksdemo get iam-policy --role eksdemo.<cluster-name>.karpenter.karpenter -o yaml
InlinePolicies:
- Name: eksctl-blue-addon-iamserviceaccount-karpenter-karpenter-Policy1
  PolicyDocument: |-
    {
        "Version": "2012-10-17",
        "Statement": [
            {
                "Action": [
                    "ec2:RunInstances",
                    "ec2:CreateFleet"
                ],
                "Resource": [
                    "arn:aws:ec2:us-west-2::image/*",
                    "arn:aws:ec2:us-west-2::snapshot/*",
                    "arn:aws:ec2:us-west-2:*:security-group/*",
                    "arn:aws:ec2:us-west-2:*:subnet/*"
                ],
                "Effect": "Allow",
                "Sid": "AllowScopedEC2InstanceAccessActions"
            },
            <snip>
            {
                "Action": "eks:DescribeCluster",
                "Resource": "arn:aws:eks:us-west-2:123456789012:cluster/blue",
                "Effect": "Allow",
                "Sid": "AllowAPIServerEndpointDiscovery"
            }
        ]
    }
ManagedPolicies: []
```

Next, let's inspect the Karpenter Node IAM Role. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get iam-policy --role KarpenterNodeRole-<cluster-name>
+------------------------------------+---------+-------------------------------------------+
|                Name                |  Type   |                Description                |
+------------------------------------+---------+-------------------------------------------+
| AmazonSSMManagedInstanceCore       | AWS Mgd | The policy for Amazon EC2 Role            |
|                                    |         | to enable AWS Systems Manager             |
|                                    |         | service core functionality.               |
+------------------------------------+---------+-------------------------------------------+
| AmazonEKS_CNI_Policy               | AWS Mgd | This policy provides the Amazon VPC       |
|                                    |         | CNI Plugin (amazon-vpc-cni-k8s) the       |
|                                    |         | permissions it requires to modify         |
|                                    |         | the IP address configuration on your      |
|                                    |         | EKS worker nodes. This permission set     |
|                                    |         | allows the CNI to list, describe, and     |
|                                    |         | modify Elastic Network Interfaces on      |
|                                    |         | your behalf. More information on the      |
|                                    |         | AWS VPC CNI Plugin is available here:     |
|                                    |         | https://github.com/aws/amazon-vpc-cni-k8s |
+------------------------------------+---------+-------------------------------------------+
| AmazonEC2ContainerRegistryReadOnly | AWS Mgd | Provides read-only access to              |
|                                    |         | Amazon EC2 Container Registry             |
|                                    |         | repositories.                             |
+------------------------------------+---------+-------------------------------------------+
| AmazonEKSWorkerNodePolicy          | AWS Mgd | This policy allows Amazon EKS             |
|                                    |         | worker nodes to connect to                |
|                                    |         | Amazon EKS Clusters.                      |
+------------------------------------+---------+-------------------------------------------+
```

If you want to view the policy document details you can run the above command again adding `-o yaml`.

```
» eksdemo get iam-policy --role KarpenterNodeRole-<cluster-name> -o yaml
```

## (Optional) Inspect Karpenter SQS Queue and EventBridge Rules

The Karpenter install process creates an SQS queue with EventBridge rules to support Native Spot Termination Handling. Use the `eksdemo get sqs-queue` command to inspect the SQS queue that was created.

```
» eksdemo get sqs-queue
+-----------+----------------+----------+----------+-----------+
|    Age    |      Name      |   Type   | Messages | In Flight |
+-----------+----------------+----------+----------+-----------+
| 7 minutes | karpenter-blue | Standard |        0 |         0 |
+-----------+----------------+----------+----------+-----------+
```

To view all the attributes of the SQS queue use the `-o yaml` output option. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get sqs-queue karpenter-<cluster-name> -o yaml
- Attributes:
    ApproximateNumberOfMessages: "0"
    ApproximateNumberOfMessagesDelayed: "0"
    ApproximateNumberOfMessagesNotVisible: "0"
    CreatedTimestamp: "1723941870"
    DelaySeconds: "0"
    LastModifiedTimestamp: "1723941871"
    MaximumMessageSize: "262144"
    MessageRetentionPeriod: "300"
    Policy: '{"Version":"2008-10-17","Id":"EC2InterruptionPolicy","Statement":[{"Effect":"Allow","Principal":{"Service":["events.amazonaws.com","sqs.amazonaws.com"]},"Action":"sqs:SendMessage","Resource":"arn:aws:sqs:us-west-2:123456789012:karpenter-blue"},{"Sid":"DenyHTTP","Effect":"Deny","Principal":"*","Action":"sqs:*","Resource":"arn:aws:sqs:us-west-2:123456789012:karpenter-blue","Condition":{"Bool":{"aws:SecureTransport":"false"}}}]}'
    QueueArn: arn:aws:sqs:us-west-2:123456789012:karpenter-blue
    ReceiveMessageWaitTimeSeconds: "0"
    SqsManagedSseEnabled: "true"
    VisibilityTimeout: "30"
  Url: https://sqs.us-west-2.amazonaws.com/123456789012/karpenter-blue
```

To support Native Node Termination Handling, four EventBridge event patterns are routed the SQS queue. Use the `eksdemo get event-rule [NAME_PREFIX]` command to inspect the rules that were created.

```
» eksdemo get event-rule karpenter
+---------+-------------------------------------+----------+
| Status  |                Name                 |   Type   |
+---------+-------------------------------------+----------+
| ENABLED | karpenter-blue-InstanceStateChange  | Standard |
| ENABLED | karpenter-blue-Rebalance            | Standard |
| ENABLED | karpenter-blue-ScheduledChange      | Standard |
| ENABLED | karpenter-blue-SpotInterruption     | Standard |
+---------+-------------------------------------+----------+
```

To see more details about the event pattern filter, use the `-o yaml` output option. Let's take a look at the InstanceStateChange rule. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo get event-rule karpenter-<cluster-name>-InstanceStateChange -o yaml
- Arn: arn:aws:events:us-west-2:123456789012:rule/karpenter-blue-InstanceStateChange
  Description: null
  EventBusName: default
  EventPattern: '{"detail-type":["EC2 Instance State-change Notification"],"source":["aws.ec2"]}'
  ManagedBy: null
  Name: karpenter-blue-InstanceStateChange
  RoleArn: null
  ScheduleExpression: null
  State: ENABLED
```
