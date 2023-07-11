package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	cognitoidp "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoUserPoolClient struct {
	*cognitoidp.Client
}

func NewCognitoUserPoolClient() *CognitoUserPoolClient {
	return &CognitoUserPoolClient{cognitoidp.NewFromConfig(GetConfig())}
}

// Creates a new Amazon Cognito user pool and sets the password policy for the pool.
func (c *CognitoUserPoolClient) CreateUserPool(name string) (*types.UserPoolType, error) {
	input := cognitoidp.CreateUserPoolInput{
		PoolName: aws.String(name),
	}

	result, err := c.Client.CreateUserPool(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.UserPool, err
}

// Deletes the specified Amazon Cognito user pool.
func (c *CognitoUserPoolClient) DeleteUserPool(id string) error {
	_, err := c.Client.DeleteUserPool(context.Background(), &cognitoidp.DeleteUserPoolInput{
		UserPoolId: aws.String(id),
	})

	return err
}

// Returns the configuration information and metadata of the specified user pool.
func (c *CognitoUserPoolClient) DescribeUserPool(id string) (*types.UserPoolType, error) {
	result, err := c.Client.DescribeUserPool(context.Background(), &cognitoidp.DescribeUserPoolInput{
		UserPoolId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result.UserPool, nil
}

// Lists the user pools associated with an AWS account.
func (c *CognitoUserPoolClient) ListUserPools() ([]types.UserPoolDescriptionType, error) {
	pools := []types.UserPoolDescriptionType{}
	pageNum := 0

	paginator := cognitoidp.NewListUserPoolsPaginator(c.Client, &cognitoidp.ListUserPoolsInput{
		MaxResults: 60,
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		pools = append(pools, out.UserPools...)
		pageNum++
	}

	return pools, nil
}
