package acm_certificate

import (
	"io"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type CertificatePrinter struct {
	certs []*types.CertificateDetail
}

func NewPrinter(certs []*types.CertificateDetail) *CertificatePrinter {
	return &CertificatePrinter{certs}
}

func (p *CertificatePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name", "Status", "In Use"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

	for _, c := range p.certs {
		age := durafmt.ParseShort(time.Since(aws.ToTime(c.CreatedAt)))

		var inUse string
		if len(c.InUseBy) > 0 {
			inUse = "Yes"
		} else {
			inUse = "No"
		}

		table.AppendRow([]string{
			age.String(),
			resourceId.FindString(aws.ToString(c.CertificateArn)),
			aws.ToString(c.DomainName),
			string(c.Status),
			inUse,
		})
	}

	table.Print(writer)

	return nil
}

func (p *CertificatePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.certs)
}

func (p *CertificatePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.certs)
}
