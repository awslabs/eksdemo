package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
)

type ElasticloadbalancingClient struct {
	*elasticloadbalancing.Client
}

func NewElasticloadbalancingClientv1() *ElasticloadbalancingClient {
	return &ElasticloadbalancingClient{elasticloadbalancing.NewFromConfig(GetConfig())}
}

func (c *ElasticloadbalancingClient) DeleteLoadBalancer(name string) error {
	_, err := c.Client.DeleteLoadBalancer(context.Background(), &elasticloadbalancing.DeleteLoadBalancerInput{
		LoadBalancerName: aws.String(name),
	})

	return err
}

func (c *ElasticloadbalancingClient) DescribeLoadBalancers(name string) ([]types.LoadBalancerDescription, error) {
	loadBalancers := []types.LoadBalancerDescription{}
	input := elasticloadbalancing.DescribeLoadBalancersInput{}
	pageNum := 0

	if name != "" {
		input.LoadBalancerNames = []string{name}
	}

	paginator := elasticloadbalancing.NewDescribeLoadBalancersPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		loadBalancers = append(loadBalancers, out.LoadBalancerDescriptions...)
		pageNum++
	}

	return loadBalancers, nil
}
