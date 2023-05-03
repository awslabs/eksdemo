package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithywaiter "github.com/aws/smithy-go/waiter"
	"github.com/jmespath/go-jmespath"
)

type GrafanaClient struct {
	*grafana.Client
}

func NewGrafanaClient() *GrafanaClient {
	return &GrafanaClient{grafana.NewFromConfig(GetConfig())}
}

func (c *GrafanaClient) CreateWorkspace(name string, auth []string, roleArn string) (*types.WorkspaceDescription, error) {
	result, err := c.Client.CreateWorkspace(context.Background(), &grafana.CreateWorkspaceInput{
		AccountAccessType:       types.AccountAccessTypeCurrentAccount,
		AuthenticationProviders: toAuthenticationProviderTypes(auth),
		PermissionType:          types.PermissionTypeServiceManaged,
		WorkspaceDataSources:    []types.DataSourceType{types.DataSourceTypePrometheus},
		WorkspaceName:           aws.String(name),
		WorkspaceRoleArn:        aws.String(roleArn),
	})

	if err != nil {
		return nil, err
	}

	err = NewWorkspaceActiveWaiter(c.Client).Wait(context.Background(),
		&grafana.DescribeWorkspaceInput{WorkspaceId: result.Workspace.Id},
		5*time.Minute,
	)

	return result.Workspace, err
}

func (c *GrafanaClient) DeleteWorkspace(id string) error {
	_, err := c.Client.DeleteWorkspace(context.Background(), &grafana.DeleteWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	return err
}

func (c *GrafanaClient) DescribeWorkspace(id string) (*types.WorkspaceDescription, error) {
	result, err := c.Client.DescribeWorkspace(context.Background(), &grafana.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result.Workspace, nil
}

func (c *GrafanaClient) ListWorkspaces() ([]types.WorkspaceSummary, error) {
	workspaces := []types.WorkspaceSummary{}
	pageNum := 0

	paginator := grafana.NewListWorkspacesPaginator(c.Client, &grafana.ListWorkspacesInput{})

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

func (c *GrafanaClient) UpdateWorkspaceAuthentication(id, samlMetadataUrl string) error {
	result, err := c.Client.DescribeWorkspace(context.Background(), &grafana.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	})
	if err != nil {
		return err
	}

	_, err = c.Client.UpdateWorkspaceAuthentication(context.Background(), &grafana.UpdateWorkspaceAuthenticationInput{
		AuthenticationProviders: result.Workspace.Authentication.Providers,
		SamlConfiguration: &types.SamlConfiguration{
			IdpMetadata: &types.IdpMetadataMemberUrl{
				Value: samlMetadataUrl,
			},
			AssertionAttributes: &types.AssertionAttributes{
				Role: aws.String("role"),
			},
			RoleValues: &types.RoleValues{
				Admin: []string{"admin"},
			},
		},
		WorkspaceId: aws.String(id),
	})

	return err
}

func toAuthenticationProviderTypes(s []string) []types.AuthenticationProviderTypes {
	apt := make([]types.AuthenticationProviderTypes, len(s))
	for i, v := range s {
		apt[i] = types.AuthenticationProviderTypes(v)
	}
	return apt
}

// DescribeWorkspaceAPIClient is a client that implements the DescribeWorkspace
// operation.
type DescribeWorkspaceAPIClient interface {
	DescribeWorkspace(context.Context, *grafana.DescribeWorkspaceInput, ...func(*grafana.Options)) (*grafana.DescribeWorkspaceOutput, error)
}

var _ DescribeWorkspaceAPIClient = (*grafana.Client)(nil)

type WorkspaceActiveWaiterOptions struct {

	// Set of options to modify how an operation is invoked. These apply to all
	// operations invoked for this client. Use functional options on operation call to
	// modify this list for per operation behavior.
	APIOptions []func(*middleware.Stack) error

	// MinDelay is the minimum amount of time to delay between retries. If unset,
	// WorkspaceActiveWaiter will use default minimum delay of 60 seconds. Note
	// that MinDelay must resolve to a value lesser than or equal to the MaxDelay.
	MinDelay time.Duration

	// MaxDelay is the maximum amount of time to delay between retries. If unset or set
	// to zero, WorkspaceActiveWaiter will use default max delay of 120 seconds.
	// Note that MaxDelay must resolve to value greater than or equal to the MinDelay.
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
	Retryable func(context.Context, *grafana.DescribeWorkspaceInput, *grafana.DescribeWorkspaceOutput, error) (bool, error)
}

// WorkspaceActiveWaiter defines the waiters for WorkspaceActive
type WorkspaceActiveWaiter struct {
	client DescribeWorkspaceAPIClient

	options WorkspaceActiveWaiterOptions
}

// NewWorkspaceActiveWaiter constructs a WorkspaceActiveWaiter.
func NewWorkspaceActiveWaiter(client DescribeWorkspaceAPIClient, optFns ...func(*WorkspaceActiveWaiterOptions)) *WorkspaceActiveWaiter {
	options := WorkspaceActiveWaiterOptions{}
	options.APIOptions = append(options.APIOptions, WaiterLogger{}.AddLogger)
	options.MinDelay = 2 * time.Second
	options.MaxDelay = 5 * time.Second
	options.Retryable = workspaceActiveStateRetryable

	for _, fn := range optFns {
		fn(&options)
	}
	return &WorkspaceActiveWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for WorkspaceActive waiter. The maxWaitDur
// is the maximum wait duration the waiter will wait. The maxWaitDur is required
// and must be greater than zero.
func (w *WorkspaceActiveWaiter) Wait(ctx context.Context, params *grafana.DescribeWorkspaceInput, maxWaitDur time.Duration, optFns ...func(*WorkspaceActiveWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for CertificateValidated waiter and
// returns the output of the successful operation. The maxWaitDur is the maximum
// wait duration the waiter will wait. The maxWaitDur is required and must be
// greater than zero.
func (w *WorkspaceActiveWaiter) WaitForOutput(ctx context.Context, params *grafana.DescribeWorkspaceInput, maxWaitDur time.Duration, optFns ...func(*WorkspaceActiveWaiterOptions)) (*grafana.DescribeWorkspaceOutput, error) {
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

		out, err := w.client.DescribeWorkspace(ctx, params, func(o *grafana.Options) {
			o.APIOptions = append(o.APIOptions, apiOptions...)
		})

		retryable, err := options.Retryable(ctx, params, out, err)
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
	return nil, fmt.Errorf("exceeded max wait time for WorkspaceActive waiter")
}

func workspaceActiveStateRetryable(ctx context.Context, input *grafana.DescribeWorkspaceInput, output *grafana.DescribeWorkspaceOutput, err error) (bool, error) {

	if err == nil {
		pathValue, err := jmespath.Search("workspace.status", output)
		if err != nil {
			return false, fmt.Errorf("error evaluating waiter state: %w", err)
		}

		expectedValue := "ACTIVE"
		value, ok := pathValue.(types.WorkspaceStatus)
		if !ok {
			return false, fmt.Errorf("waiter comparator expected types.WorkspaceStatus value, got %T", pathValue)
		}

		if string(value) == expectedValue {
			return false, nil
		}
	}

	return true, nil
}
