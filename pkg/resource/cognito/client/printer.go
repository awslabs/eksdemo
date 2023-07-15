package client

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	appClients []*types.UserPoolClientType
}

func NewPrinter(pools []*types.UserPoolClientType) *Printer {
	return &Printer{pools}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name"})

	for _, ac := range p.appClients {
		age := durafmt.ParseShort(time.Since(aws.ToTime(ac.CreationDate)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(ac.ClientId),
			aws.ToString(ac.ClientName),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.appClients)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.appClients)
}
