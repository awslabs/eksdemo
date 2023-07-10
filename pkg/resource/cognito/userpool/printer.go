package userpool

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	pools []*types.UserPoolType
}

func NewPrinter(pools []*types.UserPoolType) *Printer {
	return &Printer{pools}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name"})

	for _, a := range p.pools {
		updated := durafmt.ParseShort(time.Since(aws.ToTime(a.CreationDate)))

		table.AppendRow([]string{
			updated.String(),
			aws.ToString(a.Id),
			aws.ToString(a.Name),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.pools)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.pools)
}
