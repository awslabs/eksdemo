package nat_gateway

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

const maxNameLength = 35

type NatGatewayPrinter struct {
	nats []types.NatGateway
}

func NewPrinter(nats []types.NatGateway) *NatGatewayPrinter {
	return &NatGatewayPrinter{nats}
}

func (p *NatGatewayPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name"})

	for _, n := range p.nats {
		age := durafmt.ParseShort(time.Since(aws.ToTime(n.CreateTime)))

		table.AppendRow([]string{
			age.String(),
			string(n.State),
			aws.ToString(n.NatGatewayId),
			p.getNatGatewayName(n),
		})
	}
	table.Print(writer)

	return nil
}

func (p *NatGatewayPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nats)
}

func (p *NatGatewayPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nats)
}

func (p *NatGatewayPrinter) getNatGatewayName(nat types.NatGateway) string {
	name := ""
	for _, tag := range nat.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)

			if len(name) > maxNameLength {
				name = name[:maxNameLength-3] + "..."
			}
			continue
		}
	}
	return name
}
