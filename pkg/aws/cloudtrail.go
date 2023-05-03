package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

type CloudtrailClient struct {
	*cloudtrail.Client
	Region string
}

func NewCloudtrailClient() *CloudtrailClient {
	config := GetConfig()

	return &CloudtrailClient{cloudtrail.NewFromConfig(config), config.Region}
}

// Returns settings information for a specified trail.
func (c *CloudtrailClient) GetTrail(trailNameOrArn string) (*types.Trail, error) {
	result, err := c.Client.GetTrail(context.Background(), &cloudtrail.GetTrailInput{
		Name: aws.String(trailNameOrArn),
	})

	if err != nil {
		return nil, err
	}

	return result.Trail, nil
}

// Returns a JSON-formatted list of information about the specified trail.
// Fields include information on delivery errors, Amazon SNS and Amazon S3 errors,
// and start and stop logging times for each trail.
func (c *CloudtrailClient) GetTrailStatus(trailNameOrArn string) (*cloudtrail.GetTrailStatusOutput, error) {
	result, err := c.Client.GetTrailStatus(context.Background(), &cloudtrail.GetTrailStatusInput{
		Name: aws.String(trailNameOrArn),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Lists trails that are in the current account.
func (c *CloudtrailClient) ListTrails() ([]types.TrailInfo, error) {
	trails := []types.TrailInfo{}
	pageNum := 0

	paginator := cloudtrail.NewListTrailsPaginator(c.Client, &cloudtrail.ListTrailsInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		trails = append(trails, out.Trails...)
		pageNum++
	}

	return trails, nil
}

// Looks up management events or CloudTrail Insights events that are captured by CloudTrail.
// You can look up events that occurred in a region within the last 90 days.
func (c *CloudtrailClient) LookupEvents(insights bool, filters []types.LookupAttribute) ([]types.Event, error) {
	events := []types.Event{}
	pageNum := 0

	input := cloudtrail.LookupEventsInput{}
	if len(filters) > 0 {
		input.LookupAttributes = filters
	}

	if insights {
		input.EventCategory = types.EventCategoryInsight
	}

	paginator := cloudtrail.NewLookupEventsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		events = append(events, out.Events...)
		pageNum++
	}

	return events, nil
}
