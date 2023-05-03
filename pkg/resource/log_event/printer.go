package log_event

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type LogEventPrinter struct {
	logEvents []types.OutputLogEvent
	timestamp bool
}

func NewPrinter(logEvents []types.OutputLogEvent, timestamp bool) *LogEventPrinter {
	return &LogEventPrinter{logEvents, timestamp}
}

func (p *LogEventPrinter) PrintTable(writer io.Writer) error {
	for _, le := range p.logEvents {
		if p.timestamp {
			age := durafmt.ParseShort(time.Since(time.Unix(aws.ToInt64(le.Timestamp)/1000, 0)))
			fmt.Printf(strings.ReplaceAll(age.InternationalString(), " ", "") + ": ")
		}
		// strings.Join and strings.Fields combination removes extra spaces
		fmt.Println(strings.Join(strings.Fields(aws.ToString(le.Message)), " "))
	}

	return nil
}

func (p *LogEventPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.logEvents)
}

func (p *LogEventPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.logEvents)
}
