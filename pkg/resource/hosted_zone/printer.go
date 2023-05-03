package hosted_zone

import (
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type HostedZonePrinter struct {
	zones []types.HostedZone
}

func NewPrinter(zones []types.HostedZone) *HostedZonePrinter {
	return &HostedZonePrinter{zones}
}

func (p *HostedZonePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Records", "Zone Id"})

	for _, z := range p.zones {
		var zoneType string
		if z.Config.PrivateZone {
			zoneType = "Private"
		} else {
			zoneType = "Public"
		}

		table.AppendRow([]string{
			strings.TrimSuffix(aws.ToString(z.Name), "."),
			zoneType,
			strconv.Itoa(int(aws.ToInt64(z.ResourceRecordSetCount))),
			strings.TrimPrefix(aws.ToString(z.Id), "/hostedzone/"),
		})
	}

	table.Print(writer)

	return nil
}

func (p *HostedZonePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.zones)
}

func (p *HostedZonePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.zones)
}
