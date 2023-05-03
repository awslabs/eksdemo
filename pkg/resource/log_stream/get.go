package log_stream

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
)

type Getter struct {
	cloudwatchlogsClient *aws.CloudwatchlogsClient
	logGroupGetter       *log_group.Getter
}

func NewGetter(cloudwatchlogsClient *aws.CloudwatchlogsClient) *Getter {
	return &Getter{cloudwatchlogsClient, log_group.NewGetter(cloudwatchlogsClient)}
}

func (g *Getter) Init() {
	if g.cloudwatchlogsClient == nil {
		g.cloudwatchlogsClient = aws.NewCloudwatchlogsClient()
	}
	g.logGroupGetter = log_group.NewGetter(g.cloudwatchlogsClient)
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	lsOptions, ok := options.(*LogStreamOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to LogStreamOptions")
	}

	logGroup, err := g.logGroupGetter.GetLogGroupByName(lsOptions.LogGroupName)
	if err != nil {
		return err
	}

	logStreams, err := g.cloudwatchlogsClient.DescribeLogStreams(name, awssdk.ToString(logGroup.LogGroupName))
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(logStreams))
}
