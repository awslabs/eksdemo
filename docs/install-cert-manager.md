# Install cert-manager

[cert-manager](https://cert-manager.io/) is a X.509 certificate controller for Kubernetes. It obtains certificates from a variety of Issuers and ensure the certificates are valid and up-to-date. The `eksdemo` install includes a `ClusterIssuer` configured to create [Let's Encrypt](https://letsencrypt.org/) certificates.

cert-manager should be used when you want to expose an `eksdemo` application using Ingress using an Ingress Controller other than the AWS Load Balancer Controller. All `eksdemo` Ingress configurations use HTTPS and cert-manager will create the certificate for the Ingress Controller to use.

1. [Prerequisites](#prerequisites)
2. [Install cert-manager](#install-cert-manager-1)
3. [(Optional) Create a certificate](#optional-create-a-certificate)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install cert-manager

This section walks through the process of installing cert-manager. The command for performing the installation is:
**`eksdemo install cert-manager -c <cluster-name>`**

Let's explore the dry run output with the `--dry-run` flag. The syntax for the command is: **`eksdemo install cert-manager -c <cluster-name> --dry-run`**. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install cert-manager -c <cluster-name> --dry-run
Creating 1 dependencies for cert-manager

Creating dependency: cert-manager-irsa

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
      name: cert-manager
      namespace: cert-manager
    roleName: eksdemo.blue.cert-manager.cert-manager
    roleOnly: true
    attachPolicy:
      Version: '2012-10-17'
      Statement:
      - Effect: Allow
        Action:
        - route53:GetChange
        Resource: arn:aws:route53:::change/*
      - Effect: Allow
        Action:
        - route53:ChangeResourceRecordSets
        - route53:ListResourceRecordSets
        Resource: arn:aws:route53:::hostedzone/*
      - Effect: Allow
        Action: route53:ListHostedZonesByName
        Resource: "*"



Helm Installer Dry Run:
+---------------------+----------------------------+
| Application Version | v1.11.0                    |
| Chart Version       | 1.11.0                     |
| Chart Repository    | https://charts.jetstack.io |
| Chart Name          | cert-manager               |
| Release Name        | cert-manager               |
| Namespace           | cert-manager               |
| Wait                | false                      |
+---------------------+----------------------------+
Set Values: []
Values File:
---
installCRDs: true
replicaCount: 1
serviceAccount:
  name: cert-manager
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.cert-manager.cert-manager
image:
  tag: v1.11.0

Creating 1 post-install resources for cert-manager
Creating post-install resource: cert-manager-cluster-issuer

Kubernetes Resource Manager Dry Run:
---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
    - dns01:
        route53:
          region: us-west-2
```

From the `--dry-run` output above, you can see there are three steps to the install:
* Create an IAM Role for Service Accounts (IRSA) for the cert-manager controller
* Install the cert-manager Helm chart
* Create a `ClusterIssuer` custom resource named "letsencrypt-prod"

The cert-manager install will be ready to immediately create certificates using Let's Encrypt. The IRSA role enables the controller to update Route 53 with records that Let's Encrypt requires for domain validation. Let's proceed with installing cert-manager. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install cert-manager -c <cluster-name>
Creating 1 dependencies for cert-manager

Creating dependency: cert-manager-irsa
2023-02-08 12:36:41 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/cluster-autoscaler) will be excluded
2023-02-08 12:36:41 [ℹ]  1 iamserviceaccount (cert-manager/cert-manager) was included (based on the include/exclude rules)
2023-02-08 12:36:41 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-02-08 12:36:41 [ℹ]  1 task: { create IAM role for serviceaccount "cert-manager/cert-manager" }
2023-02-08 12:36:41 [ℹ]  building iamserviceaccount stack "eksctl-blue-addon-iamserviceaccount-cert-manager-cert-manager"
2023-02-08 12:36:41 [ℹ]  deploying stack "eksctl-blue-addon-iamserviceaccount-cert-manager-cert-manager"
2023-02-08 12:36:42 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-cert-manager-cert-manager"
2023-02-08 12:37:12 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-cert-manager-cert-manager"
2023-02-08 12:37:44 [ℹ]  waiting for CloudFormation stack "eksctl-blue-addon-iamserviceaccount-cert-manager-cert-manager"
Downloading Chart: https://charts.jetstack.io/charts/cert-manager-v1.11.0.tgz
Helm installing...
2023/02/08 12:37:59 creating 1 resource(s)
2023/02/08 12:38:00 creating 45 resource(s)
2023/02/08 12:38:01 Starting delete for "cert-manager-startupapicheck" ServiceAccount
2023/02/08 12:38:02 serviceaccounts "cert-manager-startupapicheck" not found
2023/02/08 12:38:02 creating 1 resource(s)
2023/02/08 12:38:02 Starting delete for "cert-manager-startupapicheck:create-cert" Role
2023/02/08 12:38:03 roles.rbac.authorization.k8s.io "cert-manager-startupapicheck:create-cert" not found
2023/02/08 12:38:03 creating 1 resource(s)
2023/02/08 12:38:03 Starting delete for "cert-manager-startupapicheck:create-cert" RoleBinding
2023/02/08 12:38:04 rolebindings.rbac.authorization.k8s.io "cert-manager-startupapicheck:create-cert" not found
2023/02/08 12:38:04 creating 1 resource(s)
2023/02/08 12:38:04 Starting delete for "cert-manager-startupapicheck" Job
2023/02/08 12:38:04 jobs.batch "cert-manager-startupapicheck" not found
2023/02/08 12:38:05 creating 1 resource(s)
2023/02/08 12:38:05 Watching for changes to Job cert-manager-startupapicheck with timeout of 5m0s
2023/02/08 12:38:05 Add/Modify event for cert-manager-startupapicheck: ADDED
2023/02/08 12:38:05 cert-manager-startupapicheck: Jobs active: 0, jobs failed: 0, jobs succeeded: 0
2023/02/08 12:38:05 Add/Modify event for cert-manager-startupapicheck: MODIFIED
2023/02/08 12:38:05 cert-manager-startupapicheck: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/02/08 12:38:10 Add/Modify event for cert-manager-startupapicheck: MODIFIED
2023/02/08 12:38:10 cert-manager-startupapicheck: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/02/08 12:38:16 Add/Modify event for cert-manager-startupapicheck: MODIFIED
2023/02/08 12:38:16 cert-manager-startupapicheck: Jobs active: 1, jobs failed: 0, jobs succeeded: 0
2023/02/08 12:38:18 Add/Modify event for cert-manager-startupapicheck: MODIFIED
2023/02/08 12:38:18 Starting delete for "cert-manager-startupapicheck" ServiceAccount
2023/02/08 12:38:18 Starting delete for "cert-manager-startupapicheck:create-cert" Role
2023/02/08 12:38:18 Starting delete for "cert-manager-startupapicheck:create-cert" RoleBinding
2023/02/08 12:38:19 Starting delete for "cert-manager-startupapicheck" Job
Using chart version "v1.11.0", installed "cert-manager" version "v1.11.0" in namespace "cert-manager"
NOTES:
cert-manager v1.11.0 has been deployed successfully!

In order to begin issuing certificates, you will need to set up a ClusterIssuer
or Issuer resource (for example, by creating a 'letsencrypt-staging' issuer).

More information on the different types of issuers and how to configure them
can be found in our documentation:

https://cert-manager.io/docs/configuration/

For information on how to configure cert-manager to automatically provision
Certificates for Ingress resources, take a look at the `ingress-shim`
documentation:

https://cert-manager.io/docs/usage/ingress/
Creating 1 post-install resources for cert-manager
Creating post-install resource: cert-manager-cluster-issuer
Creating ClusterIssuer "letsencrypt-prod"
```

## (Optional) Create a certificate

To test the cert-manager install you can create a test certificate. You will need a Route 53 hosted zone configured with a domain that you own. Choose a domain you would like for your certificate. Replace `<example.com>` with you domain.

```
export TEST_DOMAIN=test.<example.com>
```

Then create a Certificate resource.

```
cat <<EOF | kubectl apply -f -
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: test
spec:
  secretName: test-cert-tls
  dnsNames:
    - $TEST_DOMAIN
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
EOF
```

You can run `kubectl get cert` to view the status of your Certificate.

```
» kubectl get cert
NAME   READY   SECRET          AGE
test   False   test-cert-tls   33s
```

It will take a few minutes for Let's Encrypt to validate the certificate. Confirm that the Route 53 records were created for validations using the **`eksdemo get dns-record -z <hosted-zone>`** command. Below we use the `records` alias for the command. Replace <example.com> with you domain.

```
» eksdemo get records -z <example.com>
+----------------------------------+-------+----------------------------------------------+
|               Name               | Type  |                    Value                     |
+----------------------------------+-------+----------------------------------------------+
| example.com                      | NS    | ns-1234.awsdns-98.co.uk.                     |
|                                  |       | ns-5678.awsdns-76.org.                       |
|                                  |       | ns-123.awsdns-45.net.                        |
|                                  |       | ns-45.awsdns-67.com.                         |
| example.com                      | SOA   | ns-1234.awsdns-98.co.uk.                     |
|                                  |       | awsdns-hostmaster.amazon.com.                |
|                                  |       | 1 7200 900 1209600 86400                     |
| _acme-challenge.test.example.com | CNAME | _354518f41374633f455edd1a64448c41.ndlxkpg... |
+----------------------------------+-------+----------------------------------------------+
```

The record starting with `_acme-challenge` is the validation record. cert-manager will delete it after the certificate is created. After a few minutes check the Certificate again.

```
» kubectl get cert
NAME   READY   SECRET          AGE
test   True    test-cert-tls   4m59s
```

Congratulations, the certificate was created!