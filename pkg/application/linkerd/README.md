# `Linkerd Service Mesh`

## Documentation

[Linkerd](https://linkerd.io/)
[Getting Started Guide](https://linkerd.io/2.15/getting-started/)
[Helm Install Guide](https://linkerd.io/2.15/tasks/install-helm/)
[Helm Chart Versions](https://linkerd.io/2.15/reference/helm-chart-version-matrix/)

## Dependencies

This application presumes that you have installed [step](https://smallstep.com/cli/) in your path.

## Install with EKSDemo

`eksdemo install linkerd linkerd-crds`
`eksdemo install linkerd linkerd-control-plane`

## Installation Arguments

n/a

## Installation Validation

The quickest way to validate Linkerd is using the [CLI](https://linkerd.io/2.15/reference/cli/):
`linkerd check`

Output should look similar to this:
```sh
kubernetes-api
--------------
√ can initialize the client
√ can query the Kubernetes API

kubernetes-version
------------------
√ is running the minimum Kubernetes API version

linkerd-existence
-----------------
√ 'linkerd-config' config map exists
√ heartbeat ServiceAccount exist
√ control plane replica sets are ready
√ no unschedulable pods
√ control plane pods are ready
√ cluster networks contains all pods
√ cluster networks contains all services

linkerd-config
--------------
√ control plane Namespace exists
√ control plane ClusterRoles exist
√ control plane ClusterRoleBindings exist
√ control plane ServiceAccounts exist
√ control plane CustomResourceDefinitions exist
√ control plane MutatingWebhookConfigurations exist
√ control plane ValidatingWebhookConfigurations exist
√ proxy-init container runs as root user if docker container runtime is used

linkerd-identity
----------------
√ certificate config is valid
√ trust anchors are using supported crypto algorithm
√ trust anchors are within their validity period
√ trust anchors are valid for at least 60 days
√ issuer cert is using supported crypto algorithm
√ issuer cert is within its validity period
√ issuer cert is valid for at least 60 days
√ issuer cert is issued by the trust anchor

linkerd-webhooks-and-apisvc-tls
-------------------------------
√ proxy-injector webhook has valid cert
√ proxy-injector cert is valid for at least 60 days
√ sp-validator webhook has valid cert
√ sp-validator cert is valid for at least 60 days
√ policy-validator webhook has valid cert
√ policy-validator cert is valid for at least 60 days

linkerd-version
---------------
√ can determine the latest version
‼ cli is up-to-date
    is running version 24.7.1 but the latest edge version is 24.7.3
    see https://linkerd.io/2/checks/#l5d-version-cli for hints

control-plane-version
---------------------
√ can retrieve the control plane version
√ control plane is up-to-date
‼ control plane and cli versions match
    control plane running edge-24.7.3 but cli running edge-24.7.1
    see https://linkerd.io/2/checks/#l5d-version-control for hints

linkerd-control-plane-proxy
---------------------------
√ control plane proxies are healthy
√ control plane proxies are up-to-date
‼ control plane proxies and cli versions match
    linkerd-destination-6d7bc6d44f-p2chs running edge-24.7.3 but cli running edge-24.7.1
    see https://linkerd.io/2/checks/#l5d-cp-proxy-cli-version for hints

linkerd-extension-checks
------------------------
√ namespace configuration for extensions

Status check results are √
```
