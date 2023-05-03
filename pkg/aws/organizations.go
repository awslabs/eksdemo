package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
)

type OrganizationsClient struct {
	*organizations.Client
}

func NewOrganizationsClient() *OrganizationsClient {
	return &OrganizationsClient{organizations.NewFromConfig(GetConfig())}
}

func (c *OrganizationsClient) CreateOrganization() (*types.Organization, error) {
	result, err := c.Client.CreateOrganization(context.Background(), &organizations.CreateOrganizationInput{})
	if err != nil {
		return nil, err
	}

	return result.Organization, nil
}

func (c *OrganizationsClient) DeleteOrganization() error {
	_, err := c.Client.DeleteOrganization(context.Background(), &organizations.DeleteOrganizationInput{})
	return err
}

func (c *OrganizationsClient) DescribeOrganization() (*types.Organization, error) {
	result, err := c.Client.DescribeOrganization(context.Background(), &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}

	return result.Organization, nil
}
