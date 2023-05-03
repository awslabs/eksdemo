package network_acl

import (
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

const maxNameLength = 30

type NetworkAclPrinter struct {
	nacls []types.NetworkAcl
}

func NewPrinter(nacls []types.NetworkAcl) *NetworkAclPrinter {
	return &NetworkAclPrinter{nacls}
}

func (p *NetworkAclPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Name", "Subnets", "Default", "VPC Id"})
	table.SetColumnAlignment([]int{
		printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_RIGHT, printer.ALIGN_LEFT, printer.ALIGN_LEFT,
	})

	for _, n := range p.nacls {
		isDefault := "No"
		if aws.ToBool(n.IsDefault) {
			isDefault = "Yes"
		}

		subnets := "-"
		if n := len(n.Associations); n > 0 {
			subnets = strconv.Itoa(n)
		}

		table.AppendRow([]string{
			aws.ToString(n.NetworkAclId),
			getName(n),
			subnets,
			isDefault,
			aws.ToString(n.VpcId),
		})

	}

	table.Print(writer)
	return nil
}

func (p *NetworkAclPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nacls)
}

func (p *NetworkAclPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nacls)
}

func getName(nacl types.NetworkAcl) string {
	name := ""
	for _, tag := range nacl.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)

			if len(name) > maxNameLength {
				name = name[:maxNameLength-3] + "..."
			}
			return name
		}
	}
	return "-"
}
