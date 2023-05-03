package security_group

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type SecurityGroupPrinter struct {
	securityGroups []types.SecurityGroup
}

func NewPrinter(securityGroups []types.SecurityGroup) *SecurityGroupPrinter {
	return &SecurityGroupPrinter{securityGroups}
}

func (p *SecurityGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Group Id", "Group Name", "Description"})

	for _, sg := range p.securityGroups {
		_ = p.getName(sg)

		table.AppendRow([]string{
			aws.ToString(sg.GroupId),
			aws.ToString(sg.GroupName),
			aws.ToString(sg.Description),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *SecurityGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.securityGroups)
}

func (p *SecurityGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.securityGroups)
}

func (p *SecurityGroupPrinter) getName(sg types.SecurityGroup) string {
	name := ""
	for _, tag := range sg.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)
			continue
		}
	}
	return name
}
