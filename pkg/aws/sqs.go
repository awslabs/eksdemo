package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	*sqs.Client
}

func NewSQSClient() *SQSClient {
	return &SQSClient{sqs.NewFromConfig(GetConfig())}
}

func (c *SQSClient) GetQueueAttributes(queueUrl string) (map[string]string, error) {
	result, err := c.Client.GetQueueAttributes(context.Background(), &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(queueUrl),
		AttributeNames: []types.QueueAttributeName{
			types.QueueAttributeNameAll,
		},
	})

	if err != nil {
		return nil, err
	}

	return result.Attributes, nil
}

func (c *SQSClient) GetQueueUrl(queueName string) (string, error) {
	result, err := c.Client.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		return "", err
	}

	return aws.ToString(result.QueueUrl), nil
}

func (c *SQSClient) ListQueues() ([]string, error) {
	queueUrls := []string{}
	pageNum := 0

	paginator := sqs.NewListQueuesPaginator(c.Client, &sqs.ListQueuesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		queueUrls = append(queueUrls, out.QueueUrls...)
		pageNum++
	}

	return queueUrls, nil
}
