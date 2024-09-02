package domain

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	domains []*sagemaker.DescribeDomainOutput
}

func NewPrinter(domains []*sagemaker.DescribeDomainOutput) *Printer {
	return &Printer{domains}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Domain"})

	for _, d := range p.domains {
		age := durafmt.ParseShort(time.Since(aws.ToTime(d.CreationTime)))

		table.AppendRow([]string{
			age.String(),
			string(d.Status),
			aws.ToString(d.DomainId),
			aws.ToString(d.DomainName),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.domains)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.domains)
}
