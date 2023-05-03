package query

import (
	"fmt"
	"os"
	"sort"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	cloudwatchlogsClient *aws.CloudwatchlogsClient
}

func NewGetter(cloudwatchlogsClient *aws.CloudwatchlogsClient) *Getter {
	return &Getter{cloudwatchlogsClient}
}

func (g *Getter) Init() {
	if g.cloudwatchlogsClient == nil {
		g.cloudwatchlogsClient = aws.NewCloudwatchlogsClient()
	}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	var query types.QueryInfo
	var queries []types.QueryInfo
	var err error

	if id != "" {
		query, err = g.GetQueryById(id)
		queries = []types.QueryInfo{query}
	} else {
		queries, err = g.GetAllQueries()
	}

	if err != nil {
		return err
	}

	// Show most recent queries at the end of the list
	sort.Slice(queries, func(i, j int) bool {
		return awssdk.ToInt64(queries[i].CreateTime) < awssdk.ToInt64(queries[j].CreateTime)
	})

	return output.Print(os.Stdout, NewPrinter(queries))
}

func (g *Getter) GetAllQueries() ([]types.QueryInfo, error) {
	return g.cloudwatchlogsClient.DescribeQueries()
}

func (g *Getter) GetQueryById(id string) (types.QueryInfo, error) {
	queries, err := g.GetAllQueries()
	if err != nil {
		return types.QueryInfo{}, err
	}

	for _, q := range queries {
		if id == awssdk.ToString(q.QueryId) {
			return q, nil
		}
	}

	return types.QueryInfo{}, fmt.Errorf("logs-insights-query %q not found", id)
}
