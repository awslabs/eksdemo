package internet_gateway

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

const maxNameLength = 30

type InternetGatewayPrinter struct {
	internetGateways []types.InternetGateway
}

func NewPrinter(internetGateways []types.InternetGateway) *InternetGatewayPrinter {
	return &InternetGatewayPrinter{internetGateways: internetGateways}
}

func (p *InternetGatewayPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "State", "Name", "VPC Id"})

	for _, ig := range p.internetGateways {
		state := "Detached"
		vpcId := "-"

		if len(ig.Attachments) > 0 {
			state = "Attached"
			vpcId = aws.ToString(ig.Attachments[0].VpcId)
		}

		table.AppendRow([]string{
			aws.ToString(ig.InternetGatewayId),
			state,
			getName(ig),
			vpcId,
		})
	}

	table.Print(writer)

	return nil
}

func (p *InternetGatewayPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.internetGateways)
}

func (p *InternetGatewayPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.internetGateways)
}

func getName(ig types.InternetGateway) string {
	name := ""
	for _, tag := range ig.Tags {
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
