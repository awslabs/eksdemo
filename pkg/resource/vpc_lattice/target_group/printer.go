package target_group

import (
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type TargetGroupPrinter struct {
	targetGroups []*vpclattice.GetTargetGroupOutput
}

func NewPrinter(targetGroups []*vpclattice.GetTargetGroupOutput) *TargetGroupPrinter {
	return &TargetGroupPrinter{targetGroups}
}

func (p *TargetGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Name", "Type", "Proto:Port"})

	for _, tg := range p.targetGroups {
		age := durafmt.ParseShort(time.Since(aws.ToTime(tg.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(tg.Status),
			aws.ToString(tg.Id),
			aws.ToString(tg.Name),
			string(tg.Type),
			string(tg.Config.Protocol) + ":" + strconv.Itoa(int(aws.ToInt32(tg.Config.Port))),
		})
	}

	table.Print(writer)

	return nil
}

func (p *TargetGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.targetGroups)
}

func (p *TargetGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.targetGroups)
}
