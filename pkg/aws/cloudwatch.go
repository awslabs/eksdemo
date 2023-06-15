package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudwatchClient struct {
	*cloudwatch.Client
}

func NewCloudwatchClient() *CloudwatchClient {
	return &CloudwatchClient{cloudwatch.NewFromConfig(GetConfig())}
}

func NewCloudwatchDimensionsFilter(dimensions []string) []types.DimensionFilter {
	filters := make([]types.DimensionFilter, 0, len(dimensions))

	for _, d := range dimensions {
		filters = append(filters, types.DimensionFilter{
			Name: aws.String(d),
		})
	}

	return filters
}

// Retrieves the specified alarms.
// You can filter the results by specifying a prefix for the alarm name, the alarm state, or a prefix for any action.
func (c *CloudwatchClient) DescribeAlarms(namePrefix string) ([]types.MetricAlarm, error) {
	alarms := []types.MetricAlarm{}
	pageNum := 0

	input := cloudwatch.DescribeAlarmsInput{}
	if namePrefix != "" {
		input.AlarmNamePrefix = aws.String(namePrefix)
	}

	paginator := cloudwatch.NewDescribeAlarmsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		alarms = append(alarms, out.MetricAlarms...)
		pageNum++
	}

	return alarms, nil
}

// List the specified metrics.
// You can use the returned metrics with GetMetricData or GetMetricStatistics to get statistical data.
func (c *CloudwatchClient) ListMetrics(dimensions []types.DimensionFilter, metricName, namespace string) ([]types.Metric, error) {
	metrics := []types.Metric{}
	pageNum := 0

	input := cloudwatch.ListMetricsInput{
		Dimensions: dimensions,
	}

	if metricName != "" {
		input.MetricName = aws.String(metricName)
	}

	if namespace != "" {
		input.Namespace = aws.String(namespace)
	}

	paginator := cloudwatch.NewListMetricsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, out.Metrics...)
		pageNum++
	}

	return metrics, nil
}
