package nodegroup

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
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
	var nodegroup *types.Nodegroup
	var nodegroups []*types.Nodegroup
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		nodegroup, err = g.GetNodeGroupByName(name, clusterName)
		nodegroups = []*types.Nodegroup{nodegroup}
	} else {
		nodegroups, err = g.GetAllNodeGroups(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nodegroups))
}

func (g *Getter) GetAllNodeGroups(clusterName string) ([]*types.Nodegroup, error) {
	nodegroupNames, err := g.eksClient.ListNodegroups(clusterName)
	nodegroups := make([]*types.Nodegroup, 0, len(nodegroupNames))

	if err != nil {
		return nil, err
	}

	for _, name := range nodegroupNames {
		result, err := g.eksClient.DescribeNodegroup(clusterName, name)
		if err != nil {
			return nil, err
		}
		nodegroups = append(nodegroups, result)
	}

	return nodegroups, nil
}

func (g *Getter) GetNodeGroupByName(name, clusterName string) (*types.Nodegroup, error) {
	nodegroup, err := g.eksClient.DescribeNodegroup(clusterName, name)

	return nodegroup, aws.FormatErrorAsMessageOnly(err)

}
