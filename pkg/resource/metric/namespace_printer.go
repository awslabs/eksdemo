package metric

import (
	"io"
	"sort"
	"strconv"

	"github.com/awslabs/eksdemo/pkg/printer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type NamespacePrinter struct {
	metrics []types.Metric
}

func NewNamespacePrinter(metrics []types.Metric) *NamespacePrinter {
	return &NamespacePrinter{metrics}
}

func (p *NamespacePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Namespace", "Metrics"})

	metricCount := map[string]int{}

	for _, m := range p.metrics {
		metricCount[aws.ToString(m.Namespace)]++
	}

	keys := make([]string, 0, len(metricCount))
	for k := range metricCount {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, namespace := range keys {
		table.AppendRow([]string{
			namespace,
			strconv.Itoa(metricCount[namespace]),
		})
	}

	table.Print(writer)

	return nil
}

func (p *NamespacePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.metrics)
}

func (p *NamespacePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.metrics)
}
