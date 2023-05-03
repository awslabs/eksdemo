package query

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

const maxLogGroupNameLength = 40

type QueryPrinter struct {
	queries []types.QueryInfo
}

func NewPrinter(queries []types.QueryInfo) *QueryPrinter {
	return &QueryPrinter{queries}
}

func (p *QueryPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.NoTextWrap()
	table.SetHeader([]string{"Age", "Status", "Id", "Log Group"})

	for _, q := range p.queries {
		age := durafmt.ParseShort(time.Since(time.Unix(aws.ToInt64(q.CreateTime)/1000, 0)))

		logGroupName := aws.ToString(q.LogGroupName)

		if numLogGroups := strings.Count(aws.ToString(q.QueryString), "SOURCE"); numLogGroups > 1 {
			logGroupName += fmt.Sprintf(" (+%d more)", numLogGroups-1)
		}

		if len(logGroupName) > maxLogGroupNameLength {
			logGroupName = logGroupName[:maxLogGroupNameLength-3] + "..."
		}

		table.AppendRow([]string{
			age.String(),
			string(q.Status),
			aws.ToString(q.QueryId),
			logGroupName,
		})
	}

	table.Print(writer)

	return nil
}

func (p *QueryPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.queries)
}

func (p *QueryPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.queries)
}
