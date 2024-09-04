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

func (c *SageMakerClient) DeleteUserProfile(domainID, userProfileName string) error {
	_, err := c.Client.DeleteUserProfile(context.Background(), &sagemaker.DeleteUserProfileInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfileName),
	})

	return err
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

func (c *SageMakerClient) DescribeUserProfile(domainID, userProfileName string) (*sagemaker.DescribeUserProfileOutput, error) {
	result, err := c.Client.DescribeUserProfile(context.Background(), &sagemaker.DescribeUserProfileInput{
		DomainId:        aws.String(domainID),
		UserProfileName: aws.String(userProfileName),
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

func (c *SageMakerClient) ListUserProfiles(domainID, userProfileNameContains string) ([]types.UserProfileDetails, error) {
	input := sagemaker.ListUserProfilesInput{}

	if domainID != "" {
		input.DomainIdEquals = aws.String(domainID)
	}

	if userProfileNameContains != "" {
		input.UserProfileNameContains = aws.String(userProfileNameContains)
	}

	result, err := c.Client.ListUserProfiles(context.Background(), &input)

	if err != nil {
		return nil, err
	}

	return result.UserProfiles, nil
}
