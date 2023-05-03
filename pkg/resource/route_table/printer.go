package route_table

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type RouteTablePrinter struct {
	routeTables []types.RouteTable
}

func NewPrinter(routeTables []types.RouteTable) *RouteTablePrinter {
	return &RouteTablePrinter{routeTables}
}

func (p *RouteTablePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Main", "Subnet(s) Associated", "VPC Id"})

	multipleSubnets := false

	for _, rt := range p.routeTables {
		main := "No"
		subnets := "-"
		subnetIds := []string{}

		for _, a := range rt.Associations {
			if aws.ToBool(a.Main) {
				main = "Yes"
			}
			if a.SubnetId != nil {
				subnetIds = append(subnetIds, aws.ToString(a.SubnetId))
			}
		}

		if len(subnetIds) > 0 {
			subnets = strings.Join(subnetIds, ", ")
		}

		if len(subnetIds) > 1 {
			multipleSubnets = true
		}

		table.AppendRow([]string{
			aws.ToString(rt.RouteTableId),
			main,
			subnets,
			aws.ToString(rt.VpcId),
		})
	}

	if multipleSubnets {
		table.SeparateRows()
	}
	table.Print(writer)

	return nil
}

func (p *RouteTablePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.routeTables)
}

func (p *RouteTablePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.routeTables)
}
