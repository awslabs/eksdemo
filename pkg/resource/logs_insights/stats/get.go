package stats

import (
	"os"

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
	results, err := g.cloudwatchlogsClient.GetQueryResults(id)
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}

	return output.Print(os.Stdout, NewPrinter(results, id))
}
