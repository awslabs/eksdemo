package iam_oidc

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type IamOidcPrinter struct {
	oidcProviders []*iam.GetOpenIDConnectProviderOutput
}

func NewPrinter(oidcProviders []*iam.GetOpenIDConnectProviderOutput) *IamOidcPrinter {
	return &IamOidcPrinter{oidcProviders}
}

func (p *IamOidcPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Provider URL"})

	for _, oidc := range p.oidcProviders {
		age := durafmt.ParseShort(time.Since(aws.ToTime(oidc.CreateDate)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(oidc.Url),
		})
	}

	table.Print(writer)

	return nil
}

func (p *IamOidcPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.oidcProviders)
}

func (p *IamOidcPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.oidcProviders)
}
