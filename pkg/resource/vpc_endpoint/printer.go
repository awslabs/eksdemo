package vpc_endpoint

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type VpcEndpointPrinter struct {
	vpcEndpoints []types.VpcEndpoint
}

func NewPrinter(vpcEndpoints []types.VpcEndpoint) *VpcEndpointPrinter {
	return &VpcEndpointPrinter{vpcEndpoints}
}

func (p *VpcEndpointPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Service Name", "VPC Id"})

	for _, ep := range p.vpcEndpoints {
		age := durafmt.ParseShort(time.Since(aws.ToTime(ep.CreationTimestamp)))

		table.AppendRow([]string{
			age.String(),
			string(ep.State),
			aws.ToString(ep.VpcEndpointId),
			aws.ToString(ep.ServiceName),
			aws.ToString(ep.VpcId),
		})
	}
	table.Print(writer)

	return nil
}

func (p *VpcEndpointPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.vpcEndpoints)
}

func (p *VpcEndpointPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.vpcEndpoints)
}
