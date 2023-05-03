package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type CloudformationClient struct {
	*cloudformation.Client
}

func NewCloudformationClient() *CloudformationClient {
	return &CloudformationClient{cloudformation.NewFromConfig(GetConfig())}
}

func (c *CloudformationClient) CreateStack(stackName, templateBody string, params map[string]string, caps []types.Capability) error {
	_, err := c.Client.CreateStack(context.Background(), &cloudformation.CreateStackInput{
		Capabilities: caps,
		Parameters:   toCloudformationParameters(params),
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(templateBody),
	})
	if err != nil {
		return err
	}

	waiter := cloudformation.NewStackCreateCompleteWaiter(c.Client, func(o *cloudformation.StackCreateCompleteWaiterOptions) {
		o.APIOptions = append(o.APIOptions, WaiterLogger{}.AddLogger)
		o.MinDelay = 2 * time.Second
		o.MaxDelay = 5 * time.Second
	})

	return waiter.Wait(context.Background(),
		&cloudformation.DescribeStacksInput{StackName: aws.String(stackName)},
		5*time.Minute,
	)
}

func (c *CloudformationClient) DeleteStack(stackName string) error {
	_, err := c.Client.DeleteStack(context.Background(), &cloudformation.DeleteStackInput{
		StackName: aws.String(stackName),
	})

	return err
}

func (c *CloudformationClient) DescribeStacks(stackName string) ([]types.Stack, error) {
	stacks := []types.Stack{}
	pageNum := 0

	input := cloudformation.DescribeStacksInput{}
	if stackName != "" {
		input.StackName = aws.String(stackName)
	}

	paginator := cloudformation.NewDescribeStacksPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		stacks = append(stacks, out.Stacks...)
		pageNum++
	}

	return stacks, nil
}

func toCloudformationParameters(tags map[string]string) (params []types.Parameter) {
	for k, v := range tags {
		params = append(params, types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}
	return
}
