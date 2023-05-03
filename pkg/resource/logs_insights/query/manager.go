package query

import (
	"context"
	"fmt"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/results"
	"github.com/spf13/cobra"
)

const fieldsCmd = `fields %s`

const filterLikeCmd = `
| filter %s like %s`

const filterRawCmd = `
| filter %s`

const sortCmd = `
| sort %s %s`

const limitCmd = `
| limit %d`

type Manager struct {
	DryRun               bool
	cloudwatchlogsClient *aws.CloudwatchlogsClient
	resultsGetter        *results.Getter
}

func (m *Manager) Init() {
	if m.cloudwatchlogsClient == nil {
		m.cloudwatchlogsClient = aws.NewCloudwatchlogsClient()
	}
	m.resultsGetter = results.NewGetter(m.cloudwatchlogsClient)
}

func (m *Manager) Create(opt resource.Options) error {
	options, ok := opt.(*QueryOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast Options to QueryOptions")
	}

	if options.LogGroup == "" {
		options.LogGroup = log_group.LogGroupNameForClusterName(options.Common().ClusterName)
	}

	end := time.Now()
	start := end.Add(-options.TimeRange)

	fields := ""
	if len(options.Fields) > 0 {
		fields = fmt.Sprintf(fieldsCmd, strings.Join(options.Fields, ", "))

		if options.Query != "" {
			fields += "\n| "
		}
	}

	filter := ""
	if options.FilterLike != "" {
		filter = fmt.Sprintf(filterLikeCmd, options.FilterField, options.FilterLike)
	}

	filterRaw := ""
	if options.FilterRaw != "" {
		filterRaw = fmt.Sprintf(filterRawCmd, options.FilterRaw)
	}

	queryString := fields + options.Query + filter + filterRaw +
		fmt.Sprintf(sortCmd, options.SortField, options.SortOrder) +
		fmt.Sprintf(limitCmd, options.Limit)

	if m.DryRun {
		return m.dryRun(options.LogGroup, queryString, start, end)
	}

	queryId, err := m.cloudwatchlogsClient.StartQuery(options.LogGroup, queryString, start, end)
	if err != nil {
		return err
	}

	fmt.Printf("Query Id %q running..", queryId)

	waiter := aws.NewQueryCompleteWaiter(m.cloudwatchlogsClient.Client, func(o *aws.QueryCompleteWaiterOptions) {
		o.APIOptions = append(o.APIOptions, aws.WaiterLogger{}.AddLogger)
		o.MinDelay = 2 * time.Second
		o.MaxDelay = 5 * time.Second
	})

	err = waiter.Wait(context.Background(),
		&cloudwatchlogs.DescribeQueriesInput{LogGroupName: awssdk.String(options.LogGroup)},
		queryId,
		5*time.Minute,
	)
	if err != nil {
		return err
	}
	fmt.Println("done!")

	return m.resultsGetter.Get(queryId, printer.Table, &results.ResultsFieldOptions{
		Field:     "@message",
		ShowStats: true,
	})
}

func (m *Manager) Delete(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(logGroupName, queryString string, start, end time.Time) error {
	fmt.Printf("\nLogs Insights Resource Manager Dry Run:\n")
	fmt.Printf("CloudWatch Logs API Call %q with request parameters:\n", "StartQuery")
	fmt.Printf("LogGroupName: %s\n", logGroupName)
	fmt.Printf("QueryString: %s\n", queryString)
	fmt.Printf("StartTime: %s\n", start.Format(time.RFC1123))
	fmt.Printf("EndTime: %s\n", end.Format(time.RFC1123))

	return nil
}
