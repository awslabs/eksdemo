package alarm

import (
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	cloudwatchClient *aws.CloudwatchClient
}

func NewGetter(cloudwatchClient *aws.CloudwatchClient) *Getter {
	return &Getter{cloudwatchClient}
}

func (g *Getter) Init() {
	if g.cloudwatchClient == nil {
		g.cloudwatchClient = aws.NewCloudwatchClient()
	}
}

func (g *Getter) Get(namePrefix string, output printer.Output, options resource.Options) error {
	alarms, err := g.cloudwatchClient.DescribeAlarms(namePrefix)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewNamespacePrinter(alarms))
}
