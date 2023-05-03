package util

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServiceAccountToken(cluster *types.Cluster, namespace, serviceAccount string) error {
	kubeContext, err := kubernetes.KubeContextForCluster(cluster)
	if err != nil {
		return err
	}

	client, err := kubernetes.Client(kubeContext)
	if err != nil {
		return err
	}

	sa, err := client.CoreV1().ServiceAccounts(namespace).Get(
		context.Background(),
		serviceAccount,
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	if len(sa.Secrets) == 0 {
		return fmt.Errorf("service account is missing token")
	}

	secret, err := client.CoreV1().Secrets(namespace).Get(
		context.Background(),
		sa.Secrets[0].Name,
		metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	token, ok := secret.Data["token"]
	if !ok {
		return fmt.Errorf("service account is missing token")
	}

	fmt.Println(string(token))

	return nil
}
