package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
)

type AMPClient struct {
	*amp.Client
}

func NewAMPClient() *AMPClient {
	return &AMPClient{amp.NewFromConfig(GetConfig())}
}

// The CreateWorkspace operation creates a workspace. A workspace is a logical space
// dedicated to the storage and querying of Prometheus metrics.
// You can have one or more workspaces in each Region in your account.
func (c *AMPClient) CreateWorkspace(alias string) (*amp.CreateWorkspaceOutput, error) {
	input := amp.CreateWorkspaceInput{}

	if alias != "" {
		input.Alias = aws.String(alias)
	}

	result, err := c.Client.CreateWorkspace(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	err = amp.NewWorkspaceActiveWaiter(c.Client).Wait(context.Background(),
		&amp.DescribeWorkspaceInput{WorkspaceId: result.WorkspaceId},
		1*time.Minute,
	)

	return result, err
}

// The DeleteWorkspace operation deletes an existing workspace.
func (c *AMPClient) DeleteWorkspace(id string) error {
	_, err := c.Client.DeleteWorkspace(context.Background(), &amp.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return err
}

// The DescribeLoggingConfiguration operation returns complete information about
// the current logging configuration of the workspace.
func (c *AMPClient) DescribeLoggingConfiguration(workspaceId string) (*types.LoggingConfigurationMetadata, error) {
	out, err := c.Client.DescribeLoggingConfiguration(context.Background(), &amp.DescribeLoggingConfigurationInput{
		WorkspaceId: aws.String(workspaceId),
	})

	if err != nil {
		return nil, err
	}

	return out.LoggingConfiguration, nil
}

// The DescribeRuleGroupsNamespace operation returns complete information about one rule groups namespace.
// To retrieve a list of rule groups namespaces, use ListRuleGroupsNamespaces.
func (c *AMPClient) DescribeRuleGroupsNamespace(name, workspaceId string) (*types.RuleGroupsNamespaceDescription, error) {
	out, err := c.Client.DescribeRuleGroupsNamespace(context.Background(), &amp.DescribeRuleGroupsNamespaceInput{
		Name:        aws.String(name),
		WorkspaceId: aws.String(workspaceId),
	})

	if err != nil {
		return nil, err
	}

	return out.RuleGroupsNamespace, nil
}

// The DescribeWorkspace operation displays information about an existing workspace.
func (c *AMPClient) DescribeWorkspace(workspaceId string) (*types.WorkspaceDescription, error) {
	out, err := c.Client.DescribeWorkspace(context.Background(), &amp.DescribeWorkspaceInput{
		WorkspaceId: aws.String(workspaceId),
	})

	if err != nil {
		return nil, err
	}

	return out.Workspace, nil
}

// The ListRuleGroupsNamespaces operation returns a list of rule groups namespaces in a workspace.
func (c *AMPClient) ListRuleGroupsNamespaces(workspaceId string) ([]types.RuleGroupsNamespaceSummary, error) {
	rules := []types.RuleGroupsNamespaceSummary{}
	pageNum := 0

	paginator := amp.NewListRuleGroupsNamespacesPaginator(c.Client, &amp.ListRuleGroupsNamespacesInput{
		WorkspaceId: aws.String(workspaceId),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		rules = append(rules, out.RuleGroupsNamespaces...)
		pageNum++
	}

	return rules, nil
}

// The ListWorkspaces operation lists all of the Amazon Managed Service for Prometheus workspaces in your account.
//  This includes workspaces being created or deleted.
func (c *AMPClient) ListWorkspaces(alias string) ([]types.WorkspaceSummary, error) {
	workspaces := []types.WorkspaceSummary{}
	pageNum := 0

	input := amp.ListWorkspacesInput{}
	if alias != "" {
		input.Alias = aws.String(alias)
	}

	paginator := amp.NewListWorkspacesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, out.Workspaces...)
		pageNum++
	}

	return workspaces, nil
}
