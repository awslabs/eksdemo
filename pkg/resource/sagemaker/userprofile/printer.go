package userprofile

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	userProfiles []*sagemaker.DescribeUserProfileOutput
}

func NewPrinter(userProfiles []*sagemaker.DescribeUserProfileOutput) *Printer {
	return &Printer{userProfiles}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "User Profile", "Domain Id"})

	for _, up := range p.userProfiles {
		age := durafmt.ParseShort(time.Since(aws.ToTime(up.CreationTime)))

		table.AppendRow([]string{
			age.String(),
			string(up.Status),
			aws.ToString(up.UserProfileName),
			aws.ToString(up.DomainId),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.userProfiles)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.userProfiles)
}
