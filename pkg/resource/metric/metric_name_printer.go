package metric

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type MetricNamePrinter struct {
	metrics []types.Metric
}

func NewMetricNamePrinter(metrics []types.Metric) *MetricNamePrinter {
	return &MetricNamePrinter{metrics}
}

func (p *MetricNamePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Metric", "Dimension(s)", "Count"})

	metricMap := map[string]map[string]int{}
	multipleDimensions := false

	for _, m := range p.metrics {
		name := aws.ToString(m.MetricName)

		if _, ok := metricMap[name]; !ok {
			metricMap[name] = map[string]int{}
		}

		dim := []string{}
		if len(m.Dimensions) > 1 {
			multipleDimensions = true
		}

		for _, d := range m.Dimensions {
			dim = append(dim, aws.ToString(d.Name))
		}
		metricMap[name][strings.Join(dim, ", ")]++
	}

	metrics := make([]string, 0, len(metricMap))
	for k := range metricMap {
		metrics = append(metrics, k)
	}
	sort.Strings(metrics)

	for _, metric := range metrics {
		for dim, count := range metricMap[metric] {
			table.AppendRow([]string{
				metric,
				dim,
				strconv.Itoa(count),
			})
		}
	}

	if multipleDimensions {
		table.SeparateRows()
	}
	table.Print(writer)

	return nil
}

func (p *MetricNamePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.metrics)
}

func (p *MetricNamePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.metrics)
}
