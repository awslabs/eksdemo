package metric

import (
	"fmt"
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

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	metricOptions, ok := options.(*CloudwatchMetricOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to CloudwatchMetricOptions")
	}

	dimensions := aws.NewCloudwatchDimensionsFilter(metricOptions.Dimensions)

	metrics, err := g.cloudwatchClient.ListMetrics(dimensions, name, metricOptions.Namespace)
	if err != nil {
		return err
	}

	if metricOptions.Values {
		return output.Print(os.Stdout, NewDimensionValuePrinter(metrics, name == "", metricOptions.cmd))
	}

	if metricOptions.Namespace != "" {
		return output.Print(os.Stdout, NewMetricNamePrinter(metrics))
	}

	return output.Print(os.Stdout, NewNamespacePrinter(metrics))
}
