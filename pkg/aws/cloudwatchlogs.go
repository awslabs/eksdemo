package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithywaiter "github.com/aws/smithy-go/waiter"
)

type CloudwatchlogsClient struct {
	*cloudwatchlogs.Client
}

func NewCloudwatchlogsClient() *CloudwatchlogsClient {
	return &CloudwatchlogsClient{cloudwatchlogs.NewFromConfig(GetConfig())}
}

// Creates a log group with the specified name.
func (c *CloudwatchlogsClient) CreateLogGroup(name string) (*cloudwatchlogs.CreateLogGroupOutput, error) {
	return c.Client.CreateLogGroup(context.Background(), &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(name),
	})
}

// Deletes the specified log group and permanently deletes all the archived log events associated with the log group.
func (c *CloudwatchlogsClient) DeleteLogGroup(name string) error {
	_, err := c.Client.DeleteLogGroup(context.Background(), &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(name),
	})

	return err
}

// Lists the specified log groups. You can list all your log groups or filter the results by prefix.
// The results are ASCII-sorted by log group name.
func (c *CloudwatchlogsClient) DescribeLogGroups(namePrefix string) ([]types.LogGroup, error) {
	logGroups := []types.LogGroup{}
	pageNum := 0

	input := cloudwatchlogs.DescribeLogGroupsInput{}
	if namePrefix != "" {
		input.LogGroupNamePrefix = aws.String(namePrefix)
	}

	paginator := cloudwatchlogs.NewDescribeLogGroupsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logGroups = append(logGroups, out.LogGroups...)
		pageNum++
	}

	return logGroups, nil
}

// Lists the log streams for the specified log group. You can list all the log streams or filter the results by prefix.
// You can also control how the results are ordered.
func (c *CloudwatchlogsClient) DescribeLogStreams(namePrefix, logGroupName string) ([]types.LogStream, error) {
	logStreams := []types.LogStream{}
	pageNum := 0

	input := cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName: aws.String(logGroupName),
	}

	if namePrefix != "" {
		input.LogStreamNamePrefix = aws.String(namePrefix)
	}

	paginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logStreams = append(logStreams, out.LogStreams...)
		pageNum++
	}

	return logStreams, nil
}

// Returns a list of CloudWatch Logs Insights queries that are scheduled, running, or have been run recently in this account.
// You can request all queries or limit it to queries of a specific log group or queries with a certain status.
func (c *CloudwatchlogsClient) DescribeQueries() ([]types.QueryInfo, error) {
	queries := []types.QueryInfo{}
	pageNum := 0

	input := cloudwatchlogs.DescribeQueriesInput{}

	paginator := NewDescribeQueriesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		queries = append(queries, out.Queries...)
		pageNum++
	}

	return queries, nil
}

// Lists log events from the specified log stream. You can list all of the log events or filter using a time range.
// By default, this operation returns as many log events as can fit in a response size of 1MB (up to 10,000 log events).
func (c *CloudwatchlogsClient) GetLogEvents(logStreamName, logGroupName string) ([]types.OutputLogEvent, error) {
	logEvents := []types.OutputLogEvent{}
	pageNum := 0

	input := cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroupName),
		LogStreamName: aws.String(logStreamName),
		StartFromHead: aws.Bool(true),
	}

	paginator := cloudwatchlogs.NewGetLogEventsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		logEvents = append(logEvents, out.Events...)
		pageNum++
	}

	return logEvents, nil
}

// Returns the results from the specified query.
// Only the fields requested in the query are returned, along with a @ptr field, which is the identifier for the log record.
// You can use the value of @ptr in a GetLogRecord operation to get the full log record.
func (c *CloudwatchlogsClient) GetQueryResults(queryId string) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	return c.Client.GetQueryResults(context.Background(), &cloudwatchlogs.GetQueryResultsInput{
		QueryId: aws.String(queryId),
	})
}

// Schedules a query of a log group using CloudWatch Logs Insights.
// You specify the log group and time range to query and the query string to use.
func (c *CloudwatchlogsClient) StartQuery(logGroupName, queryString string, start, end time.Time) (string, error) {
	input := cloudwatchlogs.StartQueryInput{
		LogGroupName: aws.String(logGroupName),
		QueryString:  aws.String(queryString),
		StartTime:    aws.Int64(start.Unix()),
		EndTime:      aws.Int64(end.Unix()),
	}

	result, err := c.Client.StartQuery(context.Background(), &input)
	if err != nil {
		return "", err
	}

	return aws.ToString(result.QueryId), nil
}

// DescribeQueriesAPIClient is a client that implements the DescribeQueries
// operation.
type DescribeQueriesAPIClient interface {
	DescribeQueries(context.Context, *cloudwatchlogs.DescribeQueriesInput, ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeQueriesOutput, error)
}

var _ DescribeQueriesAPIClient = (*cloudwatchlogs.Client)(nil)

// DescribeQueriesPaginatorOptions is the paginator options for DescribeQueries
type DescribeQueriesPaginatorOptions struct {
	// Limits the number of returned queries to the specified number.
	MaxResults int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// DescribeQueriesPaginator is a paginator for DescribeQueries
type DescribeQueriesPaginator struct {
	options   DescribeQueriesPaginatorOptions
	client    DescribeQueriesAPIClient
	params    *cloudwatchlogs.DescribeQueriesInput
	nextToken *string
	firstPage bool
}

// NewDescribeQueriesPaginator returns a new DescribeQueriesPaginator
func NewDescribeQueriesPaginator(client DescribeQueriesAPIClient, params *cloudwatchlogs.DescribeQueriesInput, optFns ...func(*DescribeQueriesPaginatorOptions)) *DescribeQueriesPaginator {
	if params == nil {
		params = &cloudwatchlogs.DescribeQueriesInput{}
	}

	options := DescribeQueriesPaginatorOptions{}
	if params.MaxResults != nil {
		options.MaxResults = *params.MaxResults
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &DescribeQueriesPaginator{
		options:   options,
		client:    client,
		params:    params,
		firstPage: true,
		nextToken: params.NextToken,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *DescribeQueriesPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextToken != nil && len(*p.nextToken) != 0)
}

// NextPage retrieves the next DescribeQueries page.
func (p *DescribeQueriesPaginator) NextPage(ctx context.Context, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeQueriesOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.NextToken = p.nextToken

	var maxResults *int32
	if p.options.MaxResults > 0 {
		maxResults = &p.options.MaxResults
	}
	params.MaxResults = maxResults

	result, err := p.client.DescribeQueries(ctx, &params, optFns...)
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

// QueryCompleteWaiterOptions are waiter options for QueryCompleteWaiter
type QueryCompleteWaiterOptions struct {

	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behavior.
	APIOptions []func(*middleware.Stack) error

	// MinDelay is the minimum amount of time to delay between retries. If unset,
	// QueryCompleteWaiter will use default minimum delay of 5 seconds. Note that
	// MinDelay must resolve to a value lesser than or equal to the MaxDelay.
	MinDelay time.Duration

	// MaxDelay is the maximum amount of time to delay between retries. If unset or set
	// to zero, QueryCompleteWaiter will use default max delay of 120 seconds. Note that
	// MaxDelay must resolve to value greater than or equal to the MinDelay.
	MaxDelay time.Duration

	// LogWaitAttempts is used to enable logging for waiter retry attempts
	LogWaitAttempts bool

	// Retryable is function that can be used to override the service defined
	// waiter-behavior based on operation output, or returned error. This function is
	// used by the waiter to decide if a state is retryable or a terminal state. By
	// default service-modeled logic will populate this option. This option can thus be
	// used to define a custom waiter state with fall-back to service-modeled waiter
	// state mutators.The function returns an error in case of a failure state. In case
	// of retry state, this function returns a bool value of true and nil error, while
	// in case of success it returns a bool value of false and nil error.
	Retryable func(context.Context, *cloudwatchlogs.DescribeQueriesInput, string, *cloudwatchlogs.DescribeQueriesOutput, error) (bool, error)
}

// QueryCompleteWaiter defines the waiters for QueryComplete
type QueryCompleteWaiter struct {
	client DescribeQueriesAPIClient

	options QueryCompleteWaiterOptions
}

// NewQueryCompleteWaiter constructs a QueryCompleteWaiter.
func NewQueryCompleteWaiter(client DescribeQueriesAPIClient, optFns ...func(*QueryCompleteWaiterOptions)) *QueryCompleteWaiter {
	options := QueryCompleteWaiterOptions{}
	options.MinDelay = 5 * time.Second
	options.MaxDelay = 120 * time.Second
	options.Retryable = queryCompleteStateRetryable

	for _, fn := range optFns {
		fn(&options)
	}
	return &QueryCompleteWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for QueryComplete waiter. The maxWaitDur is the
// maximum wait duration the waiter will wait. The maxWaitDur is required and must
// be greater than zero.
func (w *QueryCompleteWaiter) Wait(ctx context.Context, params *cloudwatchlogs.DescribeQueriesInput, queryId string, maxWaitDur time.Duration, optFns ...func(*QueryCompleteWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, queryId, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for QueryComplete waiter and returns the
// output of the successful operation. The maxWaitDur is the maximum wait duration
// the waiter will wait. The maxWaitDur is required and must be greater than zero.
func (w *QueryCompleteWaiter) WaitForOutput(ctx context.Context, params *cloudwatchlogs.DescribeQueriesInput, queryId string, maxWaitDur time.Duration, optFns ...func(*QueryCompleteWaiterOptions)) (*cloudwatchlogs.DescribeQueriesOutput, error) {
	if maxWaitDur <= 0 {
		return nil, fmt.Errorf("maximum wait time for waiter must be greater than zero")
	}

	options := w.options
	for _, fn := range optFns {
		fn(&options)
	}

	if options.MaxDelay <= 0 {
		options.MaxDelay = 120 * time.Second
	}

	if options.MinDelay > options.MaxDelay {
		return nil, fmt.Errorf("minimum waiter delay %v must be lesser than or equal to maximum waiter delay of %v", options.MinDelay, options.MaxDelay)
	}

	ctx, cancelFn := context.WithTimeout(ctx, maxWaitDur)
	defer cancelFn()

	logger := smithywaiter.Logger{}
	remainingTime := maxWaitDur

	var attempt int64
	for {

		attempt++
		apiOptions := options.APIOptions
		start := time.Now()

		if options.LogWaitAttempts {
			logger.Attempt = attempt
			apiOptions = append([]func(*middleware.Stack) error{}, options.APIOptions...)
			apiOptions = append(apiOptions, logger.AddLogger)
		}

		out, err := w.client.DescribeQueries(ctx, params, func(o *cloudwatchlogs.Options) {
			o.APIOptions = append(o.APIOptions, apiOptions...)
		})

		retryable, err := options.Retryable(ctx, params, queryId, out, err)
		if err != nil {
			return nil, err
		}
		if !retryable {
			return out, nil
		}

		remainingTime -= time.Since(start)
		if remainingTime < options.MinDelay || remainingTime <= 0 {
			break
		}

		// compute exponential backoff between waiter retries
		delay, err := smithywaiter.ComputeDelay(
			attempt, options.MinDelay, options.MaxDelay, remainingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error computing waiter delay, %w", err)
		}

		remainingTime -= delay
		// sleep for the delay amount before invoking a request
		if err := smithytime.SleepWithContext(ctx, delay); err != nil {
			return nil, fmt.Errorf("request cancelled while waiting, %w", err)
		}
	}
	return nil, fmt.Errorf("exceeded max wait time for QueryComplete waiter")
}

func queryCompleteStateRetryable(ctx context.Context, input *cloudwatchlogs.DescribeQueriesInput, queryId string, output *cloudwatchlogs.DescribeQueriesOutput, err error) (bool, error) {
	if err != nil {
		return true, nil
	}

	for _, q := range output.Queries {
		if aws.ToString(q.QueryId) != queryId {
			continue
		}

		switch q.Status {
		case types.QueryStatusComplete:
			return false, nil

		case types.QueryStatusScheduled:
			return true, nil

		case types.QueryStatusRunning:
			return true, nil

		default:
			return false, fmt.Errorf("QueryId %q is in %q state", aws.ToString(q.QueryId), string(q.Status))
		}

	}

	return false, fmt.Errorf("QueryId %q not found in DescribeQueries output", queryId)
}
