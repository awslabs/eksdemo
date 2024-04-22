# Install HashiCorp Consul Application with Consul self-signed TLS

`eksdemo` makes it extremely easy to install applications from it’s extensive application catalog in your EKS clusters. In this section we will walk through the installation of [HashiCorp Consul](https://www.hashicorp.com/products/consul).

1. [Prerequisites](#prerequisites)
2. [Install HashiCorp Consul](#Install-HashiCorp-Consul) — Will use optional configuration flags to specify an Ingress with TLS

## Prerequisites

Click on the name of the prerequisites in the list below for a link to detailed instructions.

* [AWS EBS CSI Driver](/docs/install-ebs-csi-driver.md) - Consul uses stateful EBS volumes for its key/value store, so this addon is required before installation.
* [AWS Load Balancer Controller](/docs/install-awslb.md) — Consul can expose the UI through a k8s service, which is accessed through a an AWS LoadBalancer. If you plan to use the Web UI when installing Consul, this addon is required.

### Install HashiCorp Consul

[HashiCorp Consul](https://www.hashicorp.com/products/consul) is a popular Service Networking/Mesh product. 

In this section we will walk through the process of installing HashiCorp Consul. The command for performing the installation is **`eksdemo install consul -c <cluster-name> --namespace consul --enableUI --replicas <1 or 3>`**

Let’s learn a bit more about the command and it’s options before we continue by using the `-h` help shorthand flag.

```
» eksdemo install consul -h
Install consul

Usage:
  eksdemo install consul [flags]

Flags:
      --chart-version string     chart version (default "1.4.1")
  -c, --cluster string           cluster to install application (required)
      --dry-run                  don't install, just print out all installation steps
      --enableUI                 Enable Consul UI
  -h, --help                     help for consul
  -n, --namespace string         namespace to install
      --replicas int             1 or 3 replicas (default 1)
      --service-account string   service account name
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous             use previous working chart/app versions ("1.4.0"/"v1.18.0")
  -v, --version string           application version (default "v1.18.1")

Global Flags:
      --profile string   use the specific profile from your credential file
      --region string    the region to use, overrides config/env settings
```

You’ll notice above there is an optional `--replicas` flag with a default of 1. If you decide to run 3 replicas, the EKS cluster you provisioned must have at least 3 nodes. Consul servers have a podAntiAffinity configured on the stateful set using `topologyKey: kubernetes.io/hostname`. 

Lastly if you want to expose the Consul UI through a LoadBalancer, the `--enableUI` flag must be set. This will create an external facing NLB with cross-AZ enabled on port 443/TCP. Consul is configured to deploy using a self-signed certificate. When you access the NLB endpoint your browser will complain about the certificate, use your browser bypass method to access. ie. Type `thisisunsafe` for Chrome.

If you do not expose the Consul UI through a LoadBalancer, you can port-forward to the consul-server k8s service with the following command.

`kubectl port-forward service/consul-server --namespace consul 8501:8501`

Manifest Installer Dry Run:

```
» eksdemo install consul -c istio-demo --enableUI --namespace consul --replicas 3 --dry-run

Helm Installer Dry Run:
+---------------------+-------------------------------------+
| Application Version | v1.18.1                             |
| Chart Version       | 1.4.1                               |
| Chart Repository    | https://helm.releases.hashicorp.com |
| Chart Name          | consul                              |
| Release Name        | consul                              |
| Namespace           | consul                              |
| Wait                | false                               |
+---------------------+-------------------------------------+
Set Values: []
Values File:
---
global:
  # The main enabled/disabled setting.
  # If true, servers, clients, Consul DNS and the Consul UI will be enabled.
  # Namespace
  namespace: consul
  # The prefix used for all resources created in the Helm chart.
  name: consul
  # The name of the datacenter that the agents should register as.
  datacenter: dc1
  # Enables TLS across the cluster to verify authenticity of the Consul servers and clients.
  tls:
    enabled: true
  # Enables ACLs across the cluster to secure access to data and APIs.
  acls:
  # If true, automatically manage ACL tokens and policies for all Consul components.
    manageSystemACLs: true

# Configures values that configure the Consul server cluster.
server:
  enabled: true
  # The number of server agents to run. This determines the fault tolerance of the cluster.
  replicas: 3

  # Contains values that configure the Consul UI.

ui:
  enabled: true
  # Registers a Kubernetes Service for the Consul UI as a LoadBalancer.
  service:
    type: LoadBalancer
    annotations: |
      service.beta.kubernetes.io/aws-load-balancer-type: "external"
      service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"
      service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"
      service.beta.kubernetes.io/aws-load-balancer-attributes: "load_balancing.cross_zone.enabled=true"


# Configures and installs the automatic Consul Connect sidecar injector.
connectInject:
  enabled: true

apiGateway:
  manageExternalCRDs: true
```