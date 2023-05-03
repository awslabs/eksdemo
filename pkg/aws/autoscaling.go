package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
)

type AutoscalingClient struct {
	*autoscaling.Client
}

func NewAutoscalingClient() *AutoscalingClient {
	return &AutoscalingClient{autoscaling.NewFromConfig(GetConfig())}
}

func (c *AutoscalingClient) DescribeAutoScalingGroups(name string) ([]types.AutoScalingGroup, error) {
	autoScalingGroups := []types.AutoScalingGroup{}
	input := autoscaling.DescribeAutoScalingGroupsInput{}
	pageNum := 0

	if name != "" {
		input.AutoScalingGroupNames = []string{name}
	}

	paginator := autoscaling.NewDescribeAutoScalingGroupsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		autoScalingGroups = append(autoScalingGroups, out.AutoScalingGroups...)
		pageNum++
	}

	return autoScalingGroups, nil
}
