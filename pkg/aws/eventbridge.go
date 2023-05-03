package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type EventBridgeClient struct {
	*eventbridge.Client
}

func NewEventBridgeClient() *EventBridgeClient {
	return &EventBridgeClient{eventbridge.NewFromConfig(GetConfig())}
}

func (c *EventBridgeClient) ListRules(namePrefix string) ([]types.Rule, error) {
	rules := []types.Rule{}
	input := eventbridge.ListRulesInput{}
	pageNum := 0

	if namePrefix != "" {
		input.NamePrefix = aws.String(namePrefix)
	}

	paginator := NewListRulesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		rules = append(rules, out.Rules...)
		pageNum++
	}

	return rules, nil
}

// ListRulesAPIClient is a client that implements the ListRules operation.
type ListRulesAPIClient interface {
	ListRules(context.Context, *eventbridge.ListRulesInput, ...func(*eventbridge.Options)) (*eventbridge.ListRulesOutput, error)
}

var _ ListRulesAPIClient = (*eventbridge.Client)(nil)

// ListRulesPaginatorOptions is the paginator options for ListRules
type ListRulesPaginatorOptions struct {
	// (Optional) The maximum number of results to return.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListRulesPaginator is a paginator for ListRules
type ListRulesPaginator struct {
	options   ListRulesPaginatorOptions
	client    ListRulesAPIClient
	params    *eventbridge.ListRulesInput
	nextToken *string
	firstPage bool
}

// NewListRulesPaginator returns a new ListRulesPaginator
func NewListRulesPaginator(client ListRulesAPIClient, params *eventbridge.ListRulesInput, optFns ...func(*ListRulesPaginatorOptions)) *ListRulesPaginator {
	if params == nil {
		params = &eventbridge.ListRulesInput{}
	}

	options := ListRulesPaginatorOptions{}
	if params.Limit != nil {
		options.Limit = *params.Limit
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListRulesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListRulesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next ListRules page.
func (p *ListRulesPaginator) NextPage(ctx context.Context, optFns ...func(*eventbridge.Options)) (*eventbridge.ListRulesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.Limit = limit

	result, err := p.client.ListRules(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevToken := p.nextToken
	p.nextToken = result.NextToken

	if p.options.StopOnDuplicateToken &&
		prevToken != nil &&
		p.nextToken != nil &&
		*prevToken == *p.nextToken {
		p.nextToken = nil
	}

	return result, nil
}
