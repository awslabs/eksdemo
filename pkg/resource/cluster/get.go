package cluster

import (
	"errors"
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type Getter struct {
	eksClient *aws.EKSClient
}

func NewGetter(eksClient *aws.EKSClient) *Getter {
	return &Getter{eksClient}
}

func (g *Getter) Init() {
	if g.eksClient == nil {
		g.eksClient = aws.NewEKSClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var cluster *types.Cluster
	var clusters []*types.Cluster
	var err error

	if name != "" {
		cluster, err = g.GetClusterByName(name)
		clusters = []*types.Cluster{cluster}
	} else {
		clusters, err = g.GetAllClusters()
	}

	if err != nil {
		return err
	}

	currentClusterUrl := kubernetes.ClusterURLForCurrentContext()

	return output.Print(os.Stdout, NewPrinter(clusters, currentClusterUrl))
}

func (g *Getter) GetAllClusters() ([]*types.Cluster, error) {
	clusterNames, err := g.eksClient.ListClusters()
	clusters := make([]*types.Cluster, 0, len(clusterNames))

	if err != nil {
		return nil, err
	}

	for _, name := range clusterNames {
		result, err := g.eksClient.DescribeCluster(name)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, result)
	}

	return clusters, nil
}

func (g *Getter) GetClusterByName(name string) (*types.Cluster, error) {
	cluster, err := g.eksClient.DescribeCluster(name)

	var rnfe *types.ResourceNotFoundException
	if err != nil && errors.As(err, &rnfe) {
		return nil, resource.NotFoundError(fmt.Sprintf("cluster %q not found", name))
	}

	return cluster, err
}
