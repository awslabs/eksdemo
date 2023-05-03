package elastic_ip

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

const maxNameLength = 30

type ElasticIpPrinter struct {
	eips []types.Address
}

func NewPrinter(eips []types.Address) *ElasticIpPrinter {
	return &ElasticIpPrinter{eips}
}

func (p *ElasticIpPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "IP Address", "Name", "Network Interface Id"})

	for _, eip := range p.eips {
		eipId := ""
		if id := aws.ToString(eip.NetworkInterfaceId); id != "" {
			eipId = id
		}

		table.AppendRow([]string{
			aws.ToString(eip.AllocationId),
			aws.ToString(eip.PublicIp),
			getName(eip),
			eipId,
		})
	}

	table.Print(writer)
	return nil
}

func (p *ElasticIpPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.eips)
}

func (p *ElasticIpPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.eips)
}

func getName(eip types.Address) string {
	name := ""
	for _, tag := range eip.Tags {
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
