package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type IAMClient struct {
	*iam.Client
}

func NewIAMClient() *IAMClient {
	return &IAMClient{iam.NewFromConfig(GetConfig())}
}

// Creates a new role for your AWS account.
func (c *IAMClient) CreateRole(assumeRolePolicy, name, path string) (*types.Role, error) {
	if path == "" {
		path = "/"
	}

	result, err := c.Client.CreateRole(context.Background(), &iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		RoleName:                 aws.String(name),
		Path:                     aws.String(path),
	})

	if err != nil {
		return nil, err
	}

	return result.Role, nil
}

// Creates an IAM role that is linked to a specific AWS service.
func (c *IAMClient) CreateServiceLinkedRole(name string) error {
	_, err := c.Client.CreateServiceLinkedRole(context.Background(), &iam.CreateServiceLinkedRoleInput{
		AWSServiceName: aws.String(name),
	})

	return err
}

// Deletes the specified role. Unlike the AWS Management Console, when you delete a role programmatically,
// you must delete the items attached to the role manually, or the deletion fails.
func (c *IAMClient) DeleteRole(name string) error {
	_, err := c.Client.DeleteRole(context.Background(), &iam.DeleteRoleInput{
		RoleName: aws.String(name),
	})

	return err
}

// Deletes the specified inline policy that is embedded in the specified IAM role.
func (c *IAMClient) DeleteRolePolicy(roleName, policyName string) error {
	_, err := c.Client.DeleteRolePolicy(context.Background(), &iam.DeleteRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})

	return err
}

// Removes the specified managed policy from the specified role.
func (c *IAMClient) DetachRolePolicy(roleName, policyArn string) error {
	_, err := c.Client.DetachRolePolicy(context.Background(), &iam.DetachRolePolicyInput{
		PolicyArn: aws.String(policyArn),
		RoleName:  aws.String(roleName),
	})

	return err
}

// Returns information about the specified OpenID Connect (OIDC) provider resource object in IAM.
func (c *IAMClient) GetOpenIDConnectProvider(arn string) (*iam.GetOpenIDConnectProviderOutput, error) {
	return c.Client.GetOpenIDConnectProvider(context.Background(), &iam.GetOpenIDConnectProviderInput{
		OpenIDConnectProviderArn: aws.String(arn),
	})
}

// Retrieves information about the specified managed policy, including the policy's default version
// and the total number of IAM users, groups, and roles to which the policy is attached
func (c *IAMClient) GetPolicy(arn string) (*types.Policy, error) {
	result, err := c.Client.GetPolicy(context.Background(), &iam.GetPolicyInput{
		PolicyArn: aws.String(arn),
	})

	if err != nil {
		return nil, err
	}

	return result.Policy, nil
}

// Retrieves information about the specified version of the specified managed policy,
// including the policy document.
func (c *IAMClient) GetPolicyVersion(arn, version string) (*types.PolicyVersion, error) {
	result, err := c.Client.GetPolicyVersion(context.Background(), &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(arn),
		VersionId: aws.String(version),
	})

	if err != nil {
		return nil, err
	}

	return result.PolicyVersion, nil
}

// Retrieves information about the specified role, including the role's path, GUID, ARN,
// and the role's trust policy that grants permission to assume the role.
func (c *IAMClient) GetRole(name string) (*types.Role, error) {
	result, err := c.Client.GetRole(context.Background(), &iam.GetRoleInput{
		RoleName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return result.Role, nil
}

// Retrieves the specified inline policy document that is embedded with the specified IAM role.
func (c *IAMClient) GetRolePolicy(policyName, roleName string) (string, error) {
	result, err := c.Client.GetRolePolicy(context.Background(), &iam.GetRolePolicyInput{
		PolicyName: aws.String(policyName),
		RoleName:   aws.String(roleName),
	})

	if err != nil {
		return "", err
	}

	return aws.ToString(result.PolicyDocument), nil
}

// Lists all managed policies that are attached to the specified IAM role.
func (c *IAMClient) ListAttachedRolePolicies(roleName string) ([]types.AttachedPolicy, error) {
	policies := []types.AttachedPolicy{}
	pageNum := 0

	paginator := iam.NewListAttachedRolePoliciesPaginator(c.Client, &iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		policies = append(policies, out.AttachedPolicies...)
		pageNum++
	}

	return policies, nil
}

// Lists information about the IAM OpenID Connect (OIDC) provider resource objects defined in the AWS account.
func (c *IAMClient) ListOpenIDConnectProviders() ([]types.OpenIDConnectProviderListEntry, error) {
	result, err := c.Client.ListOpenIDConnectProviders(context.Background(), &iam.ListOpenIDConnectProvidersInput{})
	if err != nil {
		return nil, err
	}

	return result.OpenIDConnectProviderList, nil
}

// Lists the names of the inline policies that are embedded in the specified IAM role.
func (c *IAMClient) ListRolePolicies(roleName string) ([]string, error) {
	policyNames := []string{}
	pageNum := 0

	paginator := iam.NewListRolePoliciesPaginator(c.Client, &iam.ListRolePoliciesInput{
		RoleName: aws.String(roleName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		policyNames = append(policyNames, out.PolicyNames...)
		pageNum++
	}

	return policyNames, nil
}

// Lists all the managed policies that are available in your AWS account,
// including your own customer-defined managed policies and all AWS managed policies.
func (c *IAMClient) ListPolicies(foo string) ([]types.Policy, error) {
	policies := []types.Policy{}
	pageNum := 0

	paginator := iam.NewListPoliciesPaginator(c.Client, &iam.ListPoliciesInput{
		MaxItems: aws.Int32(1000),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		policies = append(policies, out.Policies...)
		pageNum++
	}

	return policies, nil
}

// Lists the IAM roles that have the specified path prefix.
// If there are none, the operation returns an empty list.
func (c *IAMClient) ListRoles() ([]types.Role, error) {
	roles := []types.Role{}
	pageNum := 0

	paginator := iam.NewListRolesPaginator(c.Client, &iam.ListRolesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		roles = append(roles, out.Roles...)
		pageNum++
	}

	return roles, nil
}

// Adds or updates an inline policy document that is embedded in the specified IAM role.
func (c *IAMClient) PutRolePolicy(roleName, policyName, policyDoc string) error {
	_, err := c.Client.PutRolePolicy(context.Background(), &iam.PutRolePolicyInput{
		PolicyDocument: aws.String(policyDoc),
		PolicyName:     aws.String(policyName),
		RoleName:       aws.String(roleName),
	})

	return err
}
