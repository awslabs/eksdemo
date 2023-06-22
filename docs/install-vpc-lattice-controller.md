# Install Amazon VPC Lattice (Gateway API) Controller

The [Amazon VPC Lattice Controller](https://github.com/aws/aws-application-networking-k8s) is an implementation of the Kubernetes [Gateway API](https://gateway-api.sigs.k8s.io/) and orchestrates Amazon VPC Lattice resources using Kubernetes Custom Resources like `Gateway` and `HTTPRoute`.

`eksdemo` makes it easy to use VPC Lattice and a single command will automate the install steps for the VPC Lattice Controller: 
1. Create the VPC Lattice Controller IAM Role (IRSA)
2. Add a Security Group Rule to the EKS Cluster Security Group to allow inbound traffic from VPC Lattice
3. Install the VPC Lattice Controller Helm Chart
4. Create the VPC Lattice Controller `GatewayClass` Custom Resource

Note: The VPC Lattice Controller is also refered to as the AWS Gateway API Controller.

1. [Prerequisites](#prerequisites)
2. [Install VPC Lattice (Gateway API) Controller](#install-amazon-vpc-lattice-gateway-api-controller)
3. [Create a VPC Lattice Service Network](#create-a-vpc-lattice-service-network)
4. [Create a VPC Lattice Service](#create-a-vpc-lattice-service)
5. [Test the VPC Lattice Service](#test-the-vpc-lattice-service)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install VPC Lattice (Gateway API) Controller

This section walks through the process of installing the VPC Lattice Controller. The command for performing the installation is:
**`eksdemo install vpc-lattice-controller -c <cluster-name>`**.

Let’s expore the command and it’s options by using the -h help shorthand flag.
```
» eksdemo install vpc-lattice-controller -h
Install vpc-lattice-controller

Usage:
  eksdemo install vpc-lattice-controller [flags]

Aliases:
  vpc-lattice-controller, gateway-api-controller, vpc-lattice, lattice

Flags:
      --chart-version string     chart version (default "v0.0.12")
  -c, --cluster string           cluster to install application (required)
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for vpc-lattice-controller
  -n, --namespace string         namespace to install (default "vpc-lattice")
      --replicas int             number of replicas for the controller deployment (default 1)
      --service-account string   service account name (default "gateway-api-controller")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous             use previous working chart/app versions ("v0.0.10"/"v0.0.10")
  -v, --version string           application version (default "v0.0.12")
```

The VPC Lattice Controller specific flag is:
* `--replicas` — `eksdemo` defaults to only 1 replica for easier log viewing in a demo environment. You can use this flag to increase to the default VPC Lattice Controller Helm chart value of 2 replicas for high availability.

Next, let's review the dry run output with the `--dry-run` flag. The syntax for the command is: **`eksdemo install vpc-lattice-controller -c <cluster-name> --dry-run`**. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install vpc-lattice-controller -c <cluster-name> --dry-run
Creating 2 dependencies for vpc-lattice-controller
Creating dependency: vpc-lattice-controller-irsa

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
      name: gateway-api-controller
      namespace: vpc-lattice
    roleName: eksdemo.blue.vpc-lattice.gateway-api-controller
    roleOnly: true
    attachPolicy:
      Version: '2012-10-17'
      Statement:
      - Effect: Allow
        Action:
        - vpc-lattice:*
        - iam:CreateServiceLinkedRole
        - ec2:DescribeVpcs
        - ec2:DescribeSubnets
        Resource: "*"

Creating dependency: amazon-vpc-lattice-security-group-rule

CloudFormation Resource Manager Dry Run:
Stack name "eksdemo-blue-amazon-vpc-lattice-security-group-rule" template:

AWSTemplateFormatVersion: "2010-09-09"
Description: Allow traffic to Amazon EKS nodes from Amazon VPC Lattice
Resources:
  GatewayApiControllerVpcLatticeIngressRule:
    Type: AWS::EC2::SecurityGroupIngress
    Properties:
      Description: Allow traffic from Amazon VPC Lattice
      GroupId: sg-099fc7d2678b381e4
      IpProtocol: -1
      SourcePrefixListId: pl-0721453c7ac4ec009


Helm Installer Dry Run:
+---------------------+----------------------------------------------------------------------------------+
| Application Version | v0.0.12                                                                          |
| Chart Version       | v0.0.12                                                                          |
| Chart Repository    | oci://public.ecr.aws/aws-application-networking-k8s/aws-gateway-controller-chart |
| Chart Name          |                                                                                  |
| Release Name        | vpc-lattice-controller                                                           |
| Namespace           | vpc-lattice                                                                      |
| Wait                | false                                                                            |
+---------------------+----------------------------------------------------------------------------------+
Set Values: []
Values File:
---
fullnameOverride: gateway-api-controller
image:
  tag: v0.0.12
deployment:
  replicas: 1
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.vpc-lattice.gateway-api-controller
  name: gateway-api-controller

Creating 1 post-install resources for vpc-lattice-controller
Creating post-install resource: amazon-vpc-lattice-gateway-class

Kubernetes Resource Manager Dry Run:
---
# Create a new Gateway Class for AWS VPC lattice provider
apiVersion: gateway.networking.k8s.io/v1beta1
kind: GatewayClass
metadata:
  name: amazon-vpc-lattice
spec:
  controllerName: application-networking.k8s.aws/gateway-api-controller
```

From the `--dry-run` output above, you can see there are four steps to the install:
1. Create the VPC Lattice Controller IAM Role (IRSA)
2. Add a Security Group Rule to the EKS Cluster Security Group to allow inbound traffic from VPC Lattice
3. Install the VPC Lattice (Gateway API) Controller Helm Chart
4. Create the VPC Lattice Controller `GatewayClass` Custom Resource

After the install is completed, the VPC Lattice Controller will be ready to create VPC Lattice resources using `Gateway` and `HTTPRoute` Custom Resources. Let's proceed with installing the VPC Lattice Controller. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install vpc-lattice-controller -c <cluster-name>
Creating 2 dependencies for vpc-lattice-controller
Creating dependency: vpc-lattice-controller-irsa
2023-06-21 13:55:24 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/ebs-csi-controller-sa) will be excluded
2023-06-21 13:55:24 [ℹ]  1 iamserviceaccount (vpc-lattice/gateway-api-controller) was included (based on the include/exclude rules)
2023-06-21 13:55:24 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-06-21 13:55:24 [ℹ]  1 task: { create IAM role for serviceaccount "vpc-lattice/gateway-api-controller" }
2023-06-21 13:55:24 [ℹ]  building iamserviceaccount stack "eksctl-blue-addon-iamserviceaccount-vpc-lattice-gateway-api-controller"
2023-06-21 13:55:25 [ℹ]  deploying stack "eksctl-blue-addon-iamserviceaccount-vpc-lattice-gateway-api-controller"
2023-06-21 13:55:25 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-vpc-lattice-gateway-api-controller"
2023-06-21 13:55:56 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-vpc-lattice-gateway-api-controller"
2023-06-21 13:56:49 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-vpc-lattice-gateway-api-controller"
Creating dependency: amazon-vpc-lattice-security-group-rule
Creating CloudFormation stack "eksdemo-blue-amazon-vpc-lattice-security-group-rule" (can take 1+ min)......done
Downloading Chart: oci://public.ecr.aws/aws-application-networking-k8s/aws-gateway-controller-chart:v0.0.12
Helm installing...
2023/06/21 13:57:05 creating 3 resource(s)
2023/06/21 13:57:06 creating 1 resource(s)
2023/06/21 13:57:06 creating 1 resource(s)
2023/06/21 13:57:06 Clearing discovery cache
2023/06/21 13:57:06 beginning wait for 5 resources with timeout of 1m0s
2023/06/21 13:57:08 creating 1 resource(s)
2023/06/21 13:57:09 creating 4 resource(s)
Using chart version "v0.0.12", installed "vpc-lattice-controller" version "v0.0.12" in namespace "vpc-lattice"
NOTES:
aws-gateway-controller-chart has been installed.
This chart deploys "public.ecr.aws/aws-application-networking-k8s/aws-gateway-controller:".

Check its status by running:
  kubectl --namespace vpc-lattice get pods -l "app.kubernetes.io/instance=vpc-lattice-controller"

The controller is running in "cluster" mode.
Creating 1 post-install resources for vpc-lattice-controller
Creating post-install resource: amazon-vpc-lattice-gateway-class
Creating GatewayClass "amazon-vpc-lattice"
```

## Create a VPC Lattice Service Network

This step will follow the VPC Lattice [Getting Started](https://github.com/aws/aws-application-networking-k8s/blob/main/docs/getstarted.md) documentation.

We'll start by creating a `Gateway` custom resource called `my-hotel` that will instruct the VPC Lattice Controller to create a Service Network.

```
kubectl apply -f https://raw.githubusercontent.com/aws/aws-application-networking-k8s/main/examples/my-hotel-gateway.yaml
```

You can view the Lattice Service Network with the `eksdemo get lattice-service-network` command:
```
» eksdemo get lattice-service-network
+------------+----------------------+----------+----------+------+-----------+
|    Age     |          Id          |   Name   | Services | VPCs | Auth Type |
+------------+----------------------+----------+----------+------+-----------+
| 29 seconds | sn-0d1eb7e1dff44ff5b | my-hotel |        0 |    1 | NONE      |
+------------+----------------------+----------+----------+------+-----------+
```

## Create a VPC Lattice Service

This step continues the VPC Lattice [Getting Started](https://github.com/aws/aws-application-networking-k8s/blob/main/docs/getstarted.md) documentation.


Deploy the Kubernetes `Deployment` and `Service` for the example parking and review applications:
```
kubectl apply -f https://raw.githubusercontent.com/aws/aws-application-networking-k8s/main/examples/parking.yaml
kubectl apply -f https://raw.githubusercontent.com/aws/aws-application-networking-k8s/main/examples/review.yaml
```

You can confirm the example applications are running with `kubectl get all`.
```
» kubectl get all
NAME                           READY   STATUS    RESTARTS   AGE
pod/parking-7c4845bbf9-nkxjq   1/1     Running   0          19s
pod/parking-7c4845bbf9-x45x9   1/1     Running   0          19s
pod/review-5f598cc475-7x4zl    1/1     Running   0          18s
pod/review-5f598cc475-f6ghh    1/1     Running   0          18s

NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)   AGE
service/kubernetes   ClusterIP   10.100.0.1       <none>        443/TCP   122m
service/parking      ClusterIP   10.100.86.210    <none>        80/TCP    19s
service/review       ClusterIP   10.100.103.157   <none>        80/TCP    18s

NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/parking   2/2     2            2           19s
deployment.apps/review    2/2     2            2           18s

NAME                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/parking-7c4845bbf9   2         2         2       19s
replicaset.apps/review-5f598cc475    2         2         2       18s
```

Deploy the `HTTPRoute` custom resource named `rates` that will create a VPC Lattice Service to do path based routing as follows:
* /parking — will route traffic to the `parking` service
* /review — will route traffice to the `review` service

```
kubectl apply -f https://raw.githubusercontent.com/aws/aws-application-networking-k8s/main/examples/rate-route-path.yaml
```

It will take 20-30 seconds for the VPC Lattice Service to be created. You can view the service with the `eksdemo get lattice-service` command:
```
» eksdemo get lattice-service
+-----------+--------+-----------------------+---------------+
|    Age    | Status |          Id           |     Name      |
+-----------+--------+-----------------------+---------------+
| 2 seconds | ACTIVE | svc-0e3b248e712320f13 | rates-default |
+-----------+--------+-----------------------+---------------+
```

## Test the VPC Lattice Service

In this section we will test VPC Lattice Service. First, let's identify the DNS name of the Lattice Service. The information is available in 2 places:
* The details of the VPC Lattice Service
* The annotations of the `HTTPRoute` custom resource

You can get the DNS name from the VPC Lattice Service by using the `eksdemo get lattice-service <service-id> -o yaml` command. Or you can use the following kubectl command retrieve the DNS name from `HTTPRoute` custom resource:

```
kubectl get httproute rates -o jsonpath='{.metadata.annotations.application-networking\.k8s\.aws/lattice-assigned-domain-name}{"\n"}'
```

The DNS name will look like: `rates-default-0e3b248e712320f13.7d67968.vpc-lattice-svcs.us-west-2.on.aws`. Then start a curl Pod in your EKS cluster to test the service connectivity.

```
kubectl run curl --rm -it --image=alpine/curl -- sh
```

The command above is an interactive session inside the curl Pod and you can test the connectivity by using the curl command to test each path. Some example output below:
```
» kubectl run curl --rm -it --image=alpine/curl -- sh
If you don't see a command prompt, try pressing enter.
/ # curl rates-default-0e3b248e712320f13.7d67968.vpc-lattice-svcs.us-west-2.on.aws/parking
Requsting to Pod(parking-7c4845bbf9-nkxjq): parking handler pod
/ # curl rates-default-0e3b248e712320f13.7d67968.vpc-lattice-svcs.us-west-2.on.aws/review
Requsting to Pod(review-5f598cc475-7x4zl): review handler pod
/ # curl rates-default-0e3b248e712320f13.7d67968.vpc-lattice-svcs.us-west-2.on.aws
Not Found/ # exit
Session ended, resume using 'kubectl attach curl -c curl -i -t' command when the pod is running
pod "curl" deleted
```

The example above makes a request to the `/parking` and `/review` paths and you can see from the output that VPC Lattice is performing path based routing and directing traffic to the correct Pods.

