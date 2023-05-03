package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type ECRClient struct {
	*ecr.Client
}

func NewECRClient() *ECRClient {
	return &ECRClient{ecr.NewFromConfig(GetConfig())}
}

func (c *ECRClient) DescribeRepositories(name string) ([]types.Repository, error) {
	repositories := []types.Repository{}
	pageNum := 0

	input := ecr.DescribeRepositoriesInput{}
	if name != "" {
		input.RepositoryNames = []string{name}
	}

	paginator := ecr.NewDescribeRepositoriesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		repositories = append(repositories, out.Repositories...)
		pageNum++
	}

	return repositories, nil
}
