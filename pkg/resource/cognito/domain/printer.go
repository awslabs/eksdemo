package domain

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type Printer struct {
	domain *types.DomainDescriptionType
}

func NewPrinter(domain *types.DomainDescriptionType) *Printer {
	return &Printer{domain}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Status", "Domain", "User Pool Id"})

	// DescribeUserPoolDomain doesn't appear to return a ResourceNotFoundException
	// So we only print out the result if the domain is not nil
	if p.domain.Domain != nil {
		table.AppendRow([]string{
			string(p.domain.Status),
			aws.ToString(p.domain.Domain),
			aws.ToString(p.domain.UserPoolId),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.domain)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.domain)
}
