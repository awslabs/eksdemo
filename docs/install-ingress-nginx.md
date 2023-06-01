# Install Ingress NGINX

[ingress-nginx](https://github.com/kubernetes/ingress-nginx) is an [Ingress controller](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/) for Kubernetes using [NGINX](https://nginx.org/) as a reverse proxy and load balancer.

1. [Prerequisites](#prerequisites)
2. [Install Ingress NGINX](#install-ingress-nginx-1)
3. [(Optional) Inspect Ingress NGINX Load Balancer](#inspect-ingress-nginx-load-balancer)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install Ingress NGINX

This section walks through the process of installing Ingress NGINX. The command for performing the installation is:
**`eksdemo install ingress-nginx -c <cluster-name>`**


Let's explore the dry run output with the `--dry-run` flag. The syntax for the command is: **`eksdemo install ingress-nginx -c <cluster-name> --dry-run`**. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install ingress-nginx -c <cluster-name> --dry-run

Helm Installer Dry Run:
+---------------------+--------------------------------------------+
| Application Version | v1.8.0                                     |
| Chart Version       | 4.7.0                                      |
| Chart Repository    | https://kubernetes.github.io/ingress-nginx |
| Chart Name          | ingress-nginx                              |
| Release Name        | ingress-nginx                              |
| Namespace           | ingress-nginx                              |
| Wait                | false                                      |
+---------------------+--------------------------------------------+
Set Values: []
Values File:
---
controller:
  image:
    tag: v1.8.0
  replicaCount: 1
  service:
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
      service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
      service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing
      service.beta.kubernetes.io/aws-load-balancer-type: nlb
    externalTrafficPolicy: Local
serviceAccount:
  name: ingress-nginx
```

From the `--dry-run` output above, you can see four annotations on the service and the `externalTrafficPolicy` set to local. This follows the [AWS deployment instructions](https://kubernetes.github.io/ingress-nginx/deploy/#aws) in the Ingress NGINX install guide.

The `service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing` annotation is added to make sure the load balancer is Internet facing in case you are running the AWS Load Balancer Controller v2.5.0 with it's default webhook that would make the load balancer internal.

When you are ready to continue, proceed with installing Ingress NGINX. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install ingress-nginx -c <cluster-name>
Downloading Chart: https://github.com/kubernetes/ingress-nginx/releases/download/helm-chart-4.7.0/ingress-nginx-4.7.0.tgz
Helm installing...
2023/05/31 19:48:55 creating 1 resource(s)
2023/05/31 19:48:56 Starting delete for "ingress-nginx-admission" ServiceAccount
<snip>
2023/05/31 19:49:17 Watching for changes to Job ingress-nginx-admission-patch with timeout of 5m0s
2023/05/31 19:49:17 Add/Modify event for ingress-nginx-admission-patch: ADDED
2023/05/31 19:49:17 ingress-nginx-admission-patch: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/05/31 19:49:19 Add/Modify event for ingress-nginx-admission-patch: MODIFIED
2023/05/31 19:49:19 ingress-nginx-admission-patch: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/05/31 19:49:19 Add/Modify event for ingress-nginx-admission-patch: MODIFIED
2023/05/31 19:49:19 ingress-nginx-admission-patch: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/05/31 19:49:21 Add/Modify event for ingress-nginx-admission-patch: MODIFIED
2023/05/31 19:49:21 ingress-nginx-admission-patch: Jobs active: 0, jobs failed: 0, jobs succeeded: 0
2023/05/31 19:49:21 Add/Modify event for ingress-nginx-admission-patch: MODIFIED
2023/05/31 19:49:21 Starting delete for "ingress-nginx-admission" ServiceAccount
2023/05/31 19:49:21 Starting delete for "ingress-nginx-admission" ClusterRole
2023/05/31 19:49:21 Starting delete for "ingress-nginx-admission" ClusterRoleBinding
2023/05/31 19:49:22 Starting delete for "ingress-nginx-admission" Role
2023/05/31 19:49:22 Starting delete for "ingress-nginx-admission" RoleBinding
2023/05/31 19:49:22 Starting delete for "ingress-nginx-admission-patch" Job
Using chart version "4.7.0", installed "ingress-nginx" version "v1.8.0" in namespace "ingress-nginx"
NOTES:
The ingress-nginx controller has been installed.
It may take a few minutes for the LoadBalancer IP to be available.
You can watch the status by running 'kubectl --namespace ingress-nginx get services -o wide -w ingress-nginx-controller'

An example Ingress that makes use of the controller:
<snip>
```

## Inspect Ingress NGINX Load Balancer

The install of Ingress NGINX includes a Kubernetes Service of type `LoadBalancer` that is configured to deploy a NLB. To inspect the load balancer use the `eksdemo get load-balancer` command.

```
» eksdemo get load-balancer
+--------+--------+----------------------------------+------+-------+-----+-----+
|  Age   | State  |               Name               | Type | Stack | AZs | SGs |
+--------+--------+----------------------------------+------+-------+-----+-----+
| 1 hour | active | a932c2d30be6840c999e3db32f5a1a8c | NLB  | ipv4  |   3 |   0 |
+--------+--------+----------------------------------+------+-------+-----+-----+
* Indicates internal load balancer
```

To view the listener configuration use the `eksdemo get listenter -L <load-balancer-name>` command. Replace `<load-balancer-name>` with the name of your load balancer.

```
» eksdemo get listener -L <load-balancer-name>
+------------------+-----------+------------------------+----------------------------------+
|        Id        | Prot:Port | Default Certificate Id |          Default Action          |
+------------------+-----------+------------------------+----------------------------------+
| 6393b01308b4db1d | TCP:80    | -                      | forward to                       |
|                  |           |                        | k8s-ingressn-ingressn-7d96c6eb1d |
+------------------+-----------+------------------------+----------------------------------+
| 646991b32c3d03f3 | TCP:443   | -                      | forward to                       |
|                  |           |                        | k8s-ingressn-ingressn-4d2db6f2f0 |
+------------------+-----------+------------------------+----------------------------------+
```

The load balancer is configured to listen on both port 80 and port 443. Each port forwards to a different target group. To see more details about the target groups, use the `eksdemo get target-group` command.

```
» eksdemo get target-group
+----------------------------------+----------+------------+----------------------------------+
|               Name               |   Type   | Proto:Port |          Load Balancer           |
+----------------------------------+----------+------------+----------------------------------+
| k8s-ingressn-ingressn-4d2db6f2f0 | instance | TCP:30970  | a932c2d30be6840c999e3db32f5a1a8c |
| k8s-ingressn-ingressn-7d96c6eb1d | instance | TCP:31976  | a932c2d30be6840c999e3db32f5a1a8c |
+----------------------------------+----------+------------+----------------------------------+
```

The install of Ingress NGINX configures the NLB to use the `Instance` target type. The Kubernetes service sets up a NodePort on each worker node. From the output above we can see that port 80 on the NLB forwards to port 31976 on the EC2 instance. And port 443 on the NLB forwards to port 30970 on the instance.