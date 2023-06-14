package alarm

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type AlarmPrinter struct {
	alarms []types.MetricAlarm
}

func NewNamespacePrinter(metrics []types.MetricAlarm) *AlarmPrinter {
	return &AlarmPrinter{metrics}
}

func (p *AlarmPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Updated", "State", "Name"})

	for _, a := range p.alarms {
		updated := durafmt.ParseShort(time.Since(aws.ToTime(a.StateUpdatedTimestamp)))

		table.AppendRow([]string{
			updated.String(),
			string(a.StateValue),
			aws.ToString(a.AlarmName),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AlarmPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.alarms)
}

func (p *AlarmPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.alarms)
}
