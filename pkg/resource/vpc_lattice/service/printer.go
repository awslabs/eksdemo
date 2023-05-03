package service

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type LatticeServicePrinter struct {
	services []*vpclattice.GetServiceOutput
}

func NewPrinter(services []*vpclattice.GetServiceOutput) *LatticeServicePrinter {
	return &LatticeServicePrinter{services}
}

func (p *LatticeServicePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Name"})

	for _, s := range p.services {
		age := durafmt.ParseShort(time.Since(aws.ToTime(s.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(s.Status),
			aws.ToString(s.Id),
			aws.ToString(s.Name),
		})
	}

	table.Print(writer)

	return nil
}

func (p *LatticeServicePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.services)
}

func (p *LatticeServicePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.services)
}
