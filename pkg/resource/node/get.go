package node

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var nodes []v1.Node
	var err error

	cluster := options.Common().Cluster

	if name != "" {
		nodes, err = g.GetNodeByName(name, cluster)
	} else {
		nodes, err = g.GetAllNodes(cluster)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nodes))
}

func (g *Getter) GetNodeByName(name string, cluster *types.Cluster) ([]v1.Node, error) {
	kubeContext, err := kubernetes.KubeContextForCluster(cluster)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.Client(kubeContext)
	if err != nil {
		return nil, err
	}

	node, err := client.CoreV1().Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return []v1.Node{*node}, nil
}

func (g *Getter) GetAllNodes(cluster *types.Cluster) ([]v1.Node, error) {
	kubeContext, err := kubernetes.KubeContextForCluster(cluster)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.Client(kubeContext)
	if err != nil {
		return nil, err
	}

	nodes, err := client.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}
