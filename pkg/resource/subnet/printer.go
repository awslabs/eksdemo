package subnet

import (
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type SubnetPrinter struct {
	subnets       []types.Subnet
	multipleCidrs bool
}

func NewPrinter(subnets []types.Subnet) *SubnetPrinter {
	return &SubnetPrinter{subnets: subnets}
}

func (p *SubnetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Zone", "IPv4 CIDR", "Free", "IPv6 CIDR"})

	for _, subnet := range p.subnets {
		v6Cidrs := make([]string, 0, len(subnet.Ipv6CidrBlockAssociationSet))
		for _, cbas := range subnet.Ipv6CidrBlockAssociationSet {
			v6Cidrs = append(v6Cidrs, aws.ToString(cbas.Ipv6CidrBlock))
		}

		if len(v6Cidrs) == 0 {
			v6Cidrs = []string{"-"}
		} else {
			p.multipleCidrs = true
		}

		table.AppendRow([]string{
			aws.ToString(subnet.SubnetId),
			aws.ToString(subnet.AvailabilityZone),
			aws.ToString(subnet.CidrBlock),
			strconv.Itoa(int(aws.ToInt32(subnet.AvailableIpAddressCount))),
			strings.Join(v6Cidrs, "\n"),
		})
	}

	if p.multipleCidrs {
		table.SeparateRows()
	}

	table.Print(writer)

	return nil
}

func (p *SubnetPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.subnets)
}

func (p *SubnetPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.subnets)
}
