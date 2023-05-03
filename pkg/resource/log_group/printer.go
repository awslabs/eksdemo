package log_group

import (
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type LogGroupPrinter struct {
	logGroups []types.LogGroup
}

func NewPrinter(logGroups []types.LogGroup) *LogGroupPrinter {
	return &LogGroupPrinter{logGroups}
}

func (p *LogGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Retention"})

	for _, lg := range p.logGroups {
		age := durafmt.ParseShort(time.Since(time.Unix(aws.ToInt64(lg.CreationTime)/1000, 0)))
		retention := strconv.Itoa(int(aws.ToInt32(lg.RetentionInDays)))

		switch retention {
		case "0":
			retention = "Never expire"
		case "1":
			retention = "1 day"
		case "3":
			retention = "3 days"
		case "5":
			retention = "5 days"
		case "7":
			retention = "1 week"
		case "14":
			retention = "2 weeks"
		case "30":
			retention = "1 month"
		case "60":
			retention = "2 months"
		case "90":
			retention = "3 months"
		case "120":
			retention = "4 months"
		case "150":
			retention = "5 months"
		case "180":
			retention = "6 months"
		case "365":
			retention = "12 months"
		case "400":
			retention = "13 months"
		case "545":
			retention = "18 months"
		case "731":
			retention = "2 years"
		case "1827":
			retention = "5 years"
		case "2192":
			retention = "6 years"
		case "2557":
			retention = "7 years"
		case "2922":
			retention = "8 years"
		case "3288":
			retention = "9 years"
		case "3653":
			retention = "10 years"
		}

		table.AppendRow([]string{
			age.String(),
			aws.ToString(lg.LogGroupName),
			retention,
		})
	}
	table.Print(writer)

	return nil
}

func (p *LogGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.logGroups)
}

func (p *LogGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.logGroups)
}
