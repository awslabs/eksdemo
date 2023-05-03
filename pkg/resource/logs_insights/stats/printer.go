package stats

import (
	"fmt"
	"io"
	"math"

	"github.com/awslabs/eksdemo/pkg/printer"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type ResultsPrinter struct {
	results *cloudwatchlogs.GetQueryResultsOutput
	queryId string
}

func NewPrinter(results *cloudwatchlogs.GetQueryResultsOutput, queryId string) *ResultsPrinter {
	return &ResultsPrinter{results, queryId}
}

func (p *ResultsPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Status", "Query Id", "Scanned", "Matched", "Bytes"})

	bytes := p.results.Statistics.BytesScanned
	bytesString := ""

	for _, unit := range []string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"} {
		if math.Abs(bytes) < 1024.0 {
			bytesString = fmt.Sprintf("%3.1f %sB", bytes, unit)
			break
		}
		bytes /= 1024.0
	}

	table.AppendRow([]string{
		string(p.results.Status),
		p.queryId,
		fmt.Sprintf("%g", p.results.Statistics.RecordsScanned),
		fmt.Sprintf("%g", p.results.Statistics.RecordsMatched),
		bytesString,
	})

	table.Print(writer)

	return nil
}

func (p *ResultsPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.results)
}

func (p *ResultsPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.results)
}
