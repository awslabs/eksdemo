package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type SSMClient struct {
	*ssm.Client
}

func NewSSMClient() *SSMClient {
	return &SSMClient{ssm.NewFromConfig(GetConfig())}
}

func (c *SSMClient) DescribeInstanceInformation(instanceId string) ([]types.InstanceInformation, error) {
	filters := []types.InstanceInformationStringFilter{}
	instances := []types.InstanceInformation{}
	pageNum := 0

	if instanceId != "" {
		filters = append(filters, types.InstanceInformationStringFilter{
			Key:    aws.String("InstanceIds"),
			Values: []string{instanceId},
		})
	}

	input := &ssm.DescribeInstanceInformationInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ssm.NewDescribeInstanceInformationPaginator(c.Client, input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		instances = append(instances, out.InstanceInformationList...)
		pageNum++
	}

	return instances, nil
}

func (c *SSMClient) DescribeParameters() ([]types.ParameterMetadata, error) {
	parameters := []types.ParameterMetadata{}
	pageNum := 0

	paginator := ssm.NewDescribeParametersPaginator(c.Client, &ssm.DescribeParametersInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, out.Parameters...)
		pageNum++
	}

	return parameters, nil
}

func (c *SSMClient) DescribeSessions(id, state string) ([]types.Session, error) {
	filters := []types.SessionFilter{}
	sessions := []types.Session{}
	pageNum := 0

	input := &ssm.DescribeSessionsInput{
		State: types.SessionState(state),
	}

	if id != "" {
		filters = append(filters, types.SessionFilter{
			Key:   types.SessionFilterKeySessionId,
			Value: aws.String(id),
		})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ssm.NewDescribeSessionsPaginator(c.Client, input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, out.Sessions...)
		pageNum++
	}

	return sessions, nil
}

func (c *SSMClient) Endpoint() (aws.Endpoint, error) {
	return ssm.NewDefaultEndpointResolver().ResolveEndpoint(region, ssm.EndpointResolverOptions{})
}

func (c *SSMClient) GetParameter(name string) (*types.Parameter, error) {
	out, err := c.Client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return out.Parameter, nil
}

func (c *SSMClient) GetParametersByPath(path string) ([]types.Parameter, error) {
	parameters := []types.Parameter{}
	pageNum := 0

	input := &ssm.GetParametersByPathInput{
		Path:      aws.String(path),
		Recursive: aws.Bool(true),
	}

	paginator := ssm.NewGetParametersByPathPaginator(c.Client, input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		parameters = append(parameters, out.Parameters...)
		pageNum++
	}

	return parameters, nil
}

func (c *SSMClient) StartSession(documentName, target string, params map[string][]string) (*ssm.StartSessionOutput, error) {
	input := &ssm.StartSessionInput{
		DocumentName: aws.String(documentName),
		Target:       aws.String(target),
	}

	if len(params) > 0 {
		input.Parameters = params
	}

	return c.Client.StartSession(context.Background(), input)
}
