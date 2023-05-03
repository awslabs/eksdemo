package log_stream

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type LogStreamPrinter struct {
	logStreams []types.LogStream
}

func NewPrinter(logStreams []types.LogStream) *LogStreamPrinter {
	return &LogStreamPrinter{logStreams}
}

func (p *LogStreamPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Last Event"})

	for _, ls := range p.logStreams {
		age := durafmt.ParseShort(time.Since(time.Unix(aws.ToInt64(ls.CreationTime)/1000, 0)))
		lastEvent := durafmt.ParseShort(time.Since(time.Unix(aws.ToInt64(ls.LastEventTimestamp)/1000, 0)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(ls.LogStreamName),
			lastEvent.String(),
		})
	}
	table.Print(writer)

	return nil
}

func (p *LogStreamPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.logStreams)
}

func (p *LogStreamPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.logStreams)
}
