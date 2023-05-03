package metric

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/awslabs/eksdemo/pkg/printer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/spf13/cobra"
)

type DimensionValuePrinter struct {
	metricNameColumn bool
	metrics          []types.Metric
	cmd              *cobra.Command
}

func NewDimensionValuePrinter(metrics []types.Metric, metricNameColumn bool, cmd *cobra.Command) *DimensionValuePrinter {
	return &DimensionValuePrinter{metricNameColumn, metrics, cmd}
}

func (p *DimensionValuePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	dimensions := []string{}

	// Set the header and sort the metrics alphabetically
	if len(p.metrics) > 0 {
		header := []string{}

		// If multiple metrics included, add "Metric" to the header
		if p.metricNameColumn {
			header = append(header, "Metric")

			// Sort by Metric name and all Dimension values
			sort.Slice(p.metrics, func(i, j int) bool {
				iDimVal := ""
				for _, d := range p.metrics[i].Dimensions {
					iDimVal += aws.ToString(d.Value)
				}

				jDimVal := ""
				for _, d := range p.metrics[j].Dimensions {
					jDimVal += aws.ToString(d.Value)
				}

				return aws.ToString(p.metrics[i].MetricName)+iDimVal < aws.ToString(p.metrics[j].MetricName)+jDimVal
			})
		} else {
			// Sort by all Dimension values
			sort.Slice(p.metrics, func(i, j int) bool {
				iDimVal := ""
				for _, d := range p.metrics[i].Dimensions {
					iDimVal += aws.ToString(d.Value)
				}

				jDimVal := ""
				for _, d := range p.metrics[j].Dimensions {
					jDimVal += aws.ToString(d.Value)
				}

				return iDimVal < jDimVal
			})
		}

		for _, d := range p.metrics[0].Dimensions {
			dimensions = append(dimensions, aws.ToString(d.Name))
		}
		table.SetHeader(append(header, dimensions...))
	}

	for _, m := range p.metrics {
		row := []string{}
		dimensionNames := []string{}

		if p.metricNameColumn {
			row = append(row, aws.ToString(m.MetricName))
		}

		for _, d := range m.Dimensions {
			dimensionNames = append(dimensionNames, aws.ToString(d.Name))
			row = append(row, aws.ToString(d.Value))
		}

		if strings.Join(dimensionNames, "") != strings.Join(dimensions, "") {
			p.cmd.SilenceUsage = false

			return fmt.Errorf("%q flag requires all filtered metrics have the same dimensions. Found dimensions %q and %q",
				"values",
				strings.Join(dimensions, ", "),
				strings.Join(dimensionNames, ", "),
			)
		}

		if len(row) > 0 {
			table.AppendRow(row)
		}
	}

	table.Print(writer)

	return nil
}

func (p *DimensionValuePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.metrics)
}

func (p *DimensionValuePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.metrics)
}
