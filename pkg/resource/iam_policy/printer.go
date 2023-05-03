package iam_policy

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type IamPolicyPrinter struct {
	policies []types.Policy
}

func NewPrinter(roles []types.Policy) *IamPolicyPrinter {
	return &IamPolicyPrinter{roles}
}

func (p *IamPolicyPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name"})

	for _, p := range p.policies {
		age := durafmt.ParseShort(time.Since(aws.ToTime(p.CreateDate)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(p.PolicyName),
		})
	}

	table.Print(writer)

	return nil
}

func (p *IamPolicyPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.policies)
}

func (p *IamPolicyPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.policies)
}
