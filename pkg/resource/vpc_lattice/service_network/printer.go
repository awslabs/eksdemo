package service_network

import (
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type ServiceNetworkPrinter struct {
	serviceNetworks []*vpclattice.GetServiceNetworkOutput
}

func NewPrinter(serviceNetworks []*vpclattice.GetServiceNetworkOutput) *ServiceNetworkPrinter {
	return &ServiceNetworkPrinter{serviceNetworks}
}

func (p *ServiceNetworkPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name", "Services", "VPCs", "Auth Type"})

	for _, sn := range p.serviceNetworks {
		age := durafmt.ParseShort(time.Since(aws.ToTime(sn.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(sn.Id),
			aws.ToString(sn.Name),
			strconv.Itoa(int(aws.ToInt64(sn.NumberOfAssociatedServices))),
			strconv.Itoa(int(aws.ToInt64(sn.NumberOfAssociatedVPCs))),
			string(sn.AuthType),
		})
	}

	table.Print(writer)

	return nil
}

func (p *ServiceNetworkPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.serviceNetworks)
}

func (p *ServiceNetworkPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.serviceNetworks)
}
