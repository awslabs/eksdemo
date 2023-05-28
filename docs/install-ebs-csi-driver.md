# Install EBS CSI Driver

The [Amazon Elastic Block Store (EBS) CSI Driver](https://github.com/kubernetes-sigs/aws-ebs-csi-driver) provides a CSI interface used by Container Orchestrators to manage the lifecycle of Amazon EBS volumes.

EKS versions 1.23 and higher enable the Kubernetes in-tree to container storage interface (CSI) volume migration feature. This means the EBS CSI Driver is required to dynamically provision EBS volumes in response to a `PersistentVolumeClaim`.

The `eksdemo` install of the EBS CSI Driver includes a gp3 `StorageClass` and will set it as the default.  List of Storage Classes created:
* `gp3` (default)
* `gp3-encrypted`

You can disable this with the `--no-storageclasses` flag.

1. [Prerequisites](#prerequisites)
2. [Install EBS CSI Driver](#install-ebs-csi-driver-1)

## Prerequisites

This tutorial requires an EKS cluster with an [IAM OIDC provider configured](https://docs.aws.amazon.com/eks/latest/userguide/enable-iam-roles-for-service-accounts.html) to support IAM Roles for Service accounts (IRSA).

You can use any `eksctl` created cluster or create your cluster with `eksdemo`.

```
» eksdemo create cluster blue
```

See the [Create Cluster documentation](/docs/create-cluster.md) for configuration options.

## Install EBS CSI Driver

This section walks through the process of installing cert-manager. The command for performing the installation is:
**`eksdemo install storage-ebs-csi -c <cluster-name>`**.

Let’s expore the command and it’s options by using the -h help shorthand flag.
```
» eksdemo install storage-ebs-csi -h
Install storage-ebs-csi

Usage:
  eksdemo install storage-ebs-csi [flags]

Aliases:
  storage-ebs-csi, storage-ebscsi, storage-ebs

Flags:
      --chart-version string     chart version (default "2.19.0")
  -c, --cluster string           cluster to install application (required)
      --dry-run                  don't install, just print out all installation steps
  -h, --help                     help for storage-ebs-csi
  -n, --namespace string         namespace to install (default "kube-system")
      --no-storageclasses        don't create the gp3 StorageClasses
      --replicas int             number of replicas for the controller deployment (default 1)
      --service-account string   service account name (default "ebs-csi-controller-sa")
      --set strings              set chart values (can specify multiple or separate values with commas: key1=val1,key2=val2)
      --use-previous             use previous working chart/app versions ("2.17.1"/"v1.16.1")
  -v, --version string           application version (default "v1.19.0")
```

The EBS CSI Driver specific flags are:
* `--no-storageclasses` — This boolean flag disables the creation of the gp3 Storage Classes.
* `--replicas` — `eksdemo` defaults to only 1 replica for easier log viewing in a demo environment. You can use this flag to increase to the default EBS CSI Driver Helm chart value of 2 replicas for high availability.

Next, let's review the dry run output with the `--dry-run` flag. The syntax for the command is: **`eksdemo install storage-ebs-csi -c <cluster-name> --dry-run`**. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install storage-ebs-csi -c <cluster-name> --dry-run
Creating 1 dependencies for storage-ebs-csi

Creating dependency: ebs-csi-irsa

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
      name: ebs-csi-controller-sa
      namespace: kube-system
    roleName: eksdemo.blue.kube-system.ebs-csi-controller-sa
    roleOnly: true
    attachPolicyARNs:
    - arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy


PreInstall Dry Run:
Mark the current default StorageClass as non-default

Helm Installer Dry Run:
+---------------------+------------------------------------------------------+
| Application Version | v1.19.0                                              |
| Chart Version       | 2.19.0                                               |
| Chart Repository    | https://kubernetes-sigs.github.io/aws-ebs-csi-driver |
| Chart Name          | aws-ebs-csi-driver                                   |
| Release Name        | storage-ebs-csi                                      |
| Namespace           | kube-system                                          |
| Wait                | false                                                |
+---------------------+------------------------------------------------------+
Set Values: []
Values File:
---
image:
  tag: v1.19.0
controller:
  region: us-west-2
  replicaCount: 1
  serviceAccount:
    name: ebs-csi-controller-sa
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/eksdemo.blue.kube-system.ebs-csi-controller-sa
storageClasses:
- name: gp3
  allowVolumeExpansion: true
  annotations:
    storageclass.kubernetes.io/is-default-class: "true"
  parameters:
    csi.storage.k8s.io/fstype: ext4
    type: gp3
  volumeBindingMode: WaitForFirstConsumer
- name: gp3-encrypted
  allowVolumeExpansion: true
  parameters:
    csi.storage.k8s.io/fstype: ext4
    encrypted: "true"
    type: gp3
  volumeBindingMode: WaitForFirstConsumer
```

From the `--dry-run` output above, you can see there are two gp3 Storage Classes created using the Helm chart values under the `storageClasses` key.

Let's proceed with installing the EBS CSI Driver. If you don't want gp3 to be your default `StorageClass` add the `--no-storageclasses` flag to the command. Replace `<cluster-name>` with the name of your EKS cluster.

```
» eksdemo install storage-ebs-csi -c <cluster-name>
Creating 1 dependencies for ebs-csi
Creating dependency: ebs-csi-irsa
2023-05-28 09:42:31 [ℹ]  4 existing iamserviceaccount(s) (awslb/aws-load-balancer-controller,external-dns/external-dns,karpenter/karpenter,kube-system/ebs-csi-controller-sa) will be excluded
2023-05-28 09:42:31 [ℹ]  1 iamserviceaccount (kube-system/ebs-csi-controller-sa) was excluded (based on the include/exclude rules)
2023-05-28 09:42:31 [!]  serviceaccounts that exist in Kubernetes will be excluded, use --override-existing-serviceaccounts to override
2023-05-28 09:42:31 [ℹ]  no tasks
Checking for default StorageClass
Marking StorageClass "gp2" as non-default...done
Downloading Chart: https://github.com/kubernetes-sigs/aws-ebs-csi-driver/releases/download/helm-chart-aws-ebs-csi-driver-2.19.0/aws-ebs-csi-driver-2.19.0.tgz
Helm installing...
2023/05/28 09:42:34 creating 1 resource(s)
2023/05/28 09:42:34 creating 18 resource(s)
Using chart version "2.19.0", installed "storage-ebs-csi" version "v1.19.0" in namespace "kube-system"
NOTES:
To verify that aws-ebs-csi-driver has started, run:

    kubectl get pod -n kube-system -l "app.kubernetes.io/name=aws-ebs-csi-driver,app.kubernetes.io/instance=storage-ebs-csi"

NOTE: The [CSI Snapshotter](https://github.com/kubernetes-csi/external-snapshotter) controller and CRDs will no longer be installed as part of this chart and moving forward will be a prerequisite of using the snap shotting functionality.
```
