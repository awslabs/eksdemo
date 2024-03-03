package accessentry

import (
	"errors"
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
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

func (g *Getter) Get(arn string, output printer.Output, options resource.Options) error {
	var accessEntry *types.AccessEntry
	var accessEntries []*types.AccessEntry
	clusterName := options.Common().ClusterName
	var err error

	if arn != "" {
		accessEntry, err = g.GetAccessEntryByArn(clusterName, arn)
		accessEntries = []*types.AccessEntry{accessEntry}
	} else {
		accessEntries, err = g.GetAllAccessEntries(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(accessEntries))
}

func (g *Getter) GetAllAccessEntries(clusterName string) ([]*types.AccessEntry, error) {
	accessEntryArns, err := g.eksClient.ListAccessEntries(clusterName)
	accessEntries := make([]*types.AccessEntry, 0, len(accessEntryArns))

	if err != nil {
		return nil, err
	}

	for _, arn := range accessEntryArns {
		result, err := g.eksClient.DescribeAccessEntry(clusterName, arn)
		if err != nil {
			return nil, err
		}
		accessEntries = append(accessEntries, result)
	}

	return accessEntries, nil
}

func (g *Getter) GetAccessEntryByArn(clusterName, arn string) (*types.AccessEntry, error) {
	accessEntry, err := g.eksClient.DescribeAccessEntry(clusterName, arn)

	var rnfe *types.ResourceNotFoundException
	if err != nil && errors.As(err, &rnfe) {
		return nil, resource.NotFoundError(fmt.Sprintf("access-entry %q not found", arn))
	}

	return accessEntry, err
}
