package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

type SageMakerClient struct {
	*sagemaker.Client
}

func NewSageMakerClient() *SageMakerClient {
	return &SageMakerClient{sagemaker.NewFromConfig(GetConfig())}
}

func (c *SageMakerClient) DescribeDomain(id string) (*sagemaker.DescribeDomainOutput, error) {
	result, err := c.Client.DescribeDomain(context.Background(), &sagemaker.DescribeDomainInput{
		DomainId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *SageMakerClient) ListDomains() ([]types.DomainDetails, error) {
	result, err := c.Client.ListDomains(context.Background(), &sagemaker.ListDomainsInput{})

	if err != nil {
		return nil, err
	}

	return result.Domains, nil
}
