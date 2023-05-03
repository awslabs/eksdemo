package results

import (
	"fmt"
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
	rfOptions, ok := options.(*ResultsFieldOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ResultsFieldOptions")
	}

	results, err := g.cloudwatchlogsClient.GetQueryResults(id)
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}

	return output.Print(os.Stdout, NewPrinter(results, rfOptions.Field, rfOptions.LogStream, id, rfOptions.ShowStats))
}
