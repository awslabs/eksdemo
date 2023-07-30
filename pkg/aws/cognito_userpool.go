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

// Creates the user pool client.
// When you create a new user pool client, token revocation is automatically activated.
func (c *CognitoUserPoolClient) CreateUserPoolClient(oauthScopes, callbackUrls []string, clientName, userPoolID string) (*types.UserPoolClientType, error) {
	input := cognitoidp.CreateUserPoolClientInput{
		AllowedOAuthFlows:               []types.OAuthFlowType{types.OAuthFlowTypeCode},
		AllowedOAuthFlowsUserPoolClient: true,
		AllowedOAuthScopes:              oauthScopes,
		CallbackURLs:                    callbackUrls,
		ClientName:                      aws.String(clientName),
		GenerateSecret:                  true,
		SupportedIdentityProviders:      []string{"COGNITO"},
		UserPoolId:                      aws.String(userPoolID),
	}

	result, err := c.Client.CreateUserPoolClient(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.UserPoolClient, nil
}

// Creates a new domain for a user pool.
func (c *CognitoUserPoolClient) CreateUserPoolDomain(domain, id string) (*cognitoidp.CreateUserPoolDomainOutput, error) {
	input := cognitoidp.CreateUserPoolDomainInput{
		Domain:     aws.String(domain),
		UserPoolId: aws.String(id),
	}

	return c.Client.CreateUserPoolDomain(context.Background(), &input)
}

// Deletes the specified Amazon Cognito user pool.
func (c *CognitoUserPoolClient) DeleteUserPool(id string) error {
	_, err := c.Client.DeleteUserPool(context.Background(), &cognitoidp.DeleteUserPoolInput{
		UserPoolId: aws.String(id),
	})

	return err
}

// Allows the developer to delete the user pool client.
func (c *CognitoUserPoolClient) DeleteUserPoolClient(clientID, userPoolID string) error {
	_, err := c.Client.DeleteUserPoolClient(context.Background(), &cognitoidp.DeleteUserPoolClientInput{
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	})

	return err
}

// Deletes a domain for a user pool.
func (c *CognitoUserPoolClient) DeleteUserPoolDomain(domain, userPoolID string) error {
	_, err := c.Client.DeleteUserPoolDomain(context.Background(), &cognitoidp.DeleteUserPoolDomainInput{
		Domain:     aws.String(domain),
		UserPoolId: aws.String(userPoolID),
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

// Client method for returning the configuration information and metadata of the specified user pool app client.
func (c *CognitoUserPoolClient) DescribeUserPoolClient(clientID, userPoolID string) (*types.UserPoolClientType, error) {
	result, err := c.Client.DescribeUserPoolClient(context.Background(), &cognitoidp.DescribeUserPoolClientInput{
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	})

	if err != nil {
		return nil, err
	}

	return result.UserPoolClient, nil
}

// Gets information about a domain.
func (c *CognitoUserPoolClient) DescribeUserPoolDomain(domain string) (*types.DomainDescriptionType, error) {
	result, err := c.Client.DescribeUserPoolDomain(context.Background(), &cognitoidp.DescribeUserPoolDomainInput{
		Domain: aws.String(domain),
	})

	if err != nil {
		return nil, err
	}

	return result.DomainDescription, nil
}

// Lists the clients that have been created for the specified user pool.
func (c *CognitoUserPoolClient) ListUserPoolClients(userPoolID string) ([]types.UserPoolClientDescription, error) {
	clients := []types.UserPoolClientDescription{}
	pageNum := 0

	paginator := cognitoidp.NewListUserPoolClientsPaginator(c.Client, &cognitoidp.ListUserPoolClientsInput{
		MaxResults: 60,
		UserPoolId: aws.String(userPoolID),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		clients = append(clients, out.UserPoolClients...)
		pageNum++
	}

	return clients, nil
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
