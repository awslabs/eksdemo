package util

import (
	"context"
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/yaml"
)

func EnablePrefixAssignment(cluster *ekstypes.Cluster) error {
	kubeContext, err := kubernetes.KubeContextForCluster(cluster)
	if err != nil {
		return err
	}

	client, err := kubernetes.Client(kubeContext)
	if err != nil {
		return err
	}

	patch, _ := yaml.YAMLToJSON([]byte(prefixAssignmentPatch))

	_, err = client.AppsV1().DaemonSets("kube-system").Patch(
		context.Background(),
		"aws-node",
		types.StrategicMergePatchType,
		patch,
		metav1.PatchOptions{},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Patched daemonset \"aws-node\" with strategic merge: %sPrefix Assignment is enabled.\n",
		prefixAssignmentPatch)

	return nil
}

func EnableSecurityGroupsForPods(cluster *ekstypes.Cluster) error {
	kubeContext, err := kubernetes.KubeContextForCluster(cluster)
	if err != nil {
		return err
	}

	client, err := kubernetes.Client(kubeContext)
	if err != nil {
		return err
	}

	patch, _ := yaml.YAMLToJSON([]byte(securityGroupsForPodsPatch))

	_, err = client.AppsV1().DaemonSets("kube-system").Patch(
		context.Background(),
		"aws-node",
		types.StrategicMergePatchType,
		patch,
		metav1.PatchOptions{},
	)
	if err != nil {
		return err
	}

	fmt.Printf("Patched daemonset \"aws-node\" with strategic merge: %sSecurity Groups for Pods is enabled.\n",
		securityGroupsForPodsPatch)

	return nil
}

const prefixAssignmentPatch = `
---
spec:
  template:
    spec:
      containers:
      - name: aws-node
        env:
        - name: ENABLE_PREFIX_DELEGATION
          value: "true"
        - name: WARM_PREFIX_TARGET
          value: "1"
...
`

const securityGroupsForPodsPatch = `
---
spec:
  template:
    spec:
      containers:
      - name: aws-node
        env:
        - name: ENABLE_POD_ENI
          value: "true"
      initContainers:
      - name: aws-vpc-cni-init
        env:
        - name: DISABLE_TCP_EARLY_DEMUX
          value: "true"
...
`
