package log_group

import (
	"fmt"
	"os"

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

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var logGroup types.LogGroup
	var logGroups []types.LogGroup
	var err error

	if options.Common().ClusterName != "" {
		logGroup, err = g.GetLogGroupByName(LogGroupNameForClusterName(options.Common().ClusterName))
		logGroups = []types.LogGroup{logGroup}
	} else {
		logGroups, err = g.cloudwatchlogsClient.DescribeLogGroups(name)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logGroups))
}

func (g *Getter) GetLogGroupByName(name string) (types.LogGroup, error) {
	logGroups, err := g.cloudwatchlogsClient.DescribeLogGroups(name)
	if err != nil {
		return types.LogGroup{}, err
	}

	if len(logGroups) > 1 {
		return types.LogGroup{}, fmt.Errorf("multiple log groups found with name %q", name)
	}

	if len(logGroups) == 0 {
		return types.LogGroup{}, fmt.Errorf("log group %q not found", name)
	}

	return logGroups[0], nil
}

func LogGroupNameForClusterName(clusterName string) string {
	return fmt.Sprintf("/aws/eks/%s/cluster", clusterName)
}
