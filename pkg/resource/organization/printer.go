package organization

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/organizations/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type OrganizationPrinter struct {
	Organization *types.Organization
}

func NewPrinter(Organization *types.Organization) *OrganizationPrinter {
	return &OrganizationPrinter{Organization}
}

func (p *OrganizationPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Feature Set", "Master Account"})

	o := p.Organization

	table.AppendRow([]string{
		aws.ToString(o.Id),
		string(o.FeatureSet),
		aws.ToString(o.MasterAccountId),
	})

	table.Print(writer)

	return nil
}

func (p *OrganizationPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Organization)
}

func (p *OrganizationPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Organization)
}
