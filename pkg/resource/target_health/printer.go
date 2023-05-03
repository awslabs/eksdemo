package target_health

import (
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type TargetHealthPrinter struct {
	targets []types.TargetHealthDescription
}

func NewPrinter(targets []types.TargetHealthDescription) *TargetHealthPrinter {
	return &TargetHealthPrinter{targets}
}

func (p *TargetHealthPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"State", "Id", "Port", "Details"})

	for _, t := range p.targets {
		table.AppendRow([]string{
			string(t.TargetHealth.State),
			aws.ToString(t.Target.Id),
			strconv.Itoa(int(aws.ToInt32(t.Target.Port))),
			aws.ToString(t.TargetHealth.Description),
		})
	}
	table.Print(writer)

	return nil
}

func (p *TargetHealthPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.targets)
}

func (p *TargetHealthPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.targets)
}
