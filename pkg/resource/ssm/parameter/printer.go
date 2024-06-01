package parameter

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	params []types.Parameter
}

func NewPrinter(params []types.Parameter) *Printer {
	return &Printer{params}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Value"})

	for _, p := range p.params {
		age := durafmt.ParseShort(time.Since(aws.ToTime(p.LastModifiedDate)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(p.Name),
			aws.ToString(p.Value),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.params)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.params)
}
