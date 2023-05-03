package vpc

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type VpcPrinter struct {
	vpcs          []types.Vpc
	multipleCidrs bool
}

func NewPrinter(vpcs []types.Vpc) *VpcPrinter {
	return &VpcPrinter{vpcs: vpcs}
}

func (p *VpcPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Name", "IPv4 CIDR(s)", "IPv6 CIDR(s)"})

	for _, vpc := range p.vpcs {
		name := p.getVpcName(vpc)
		if aws.ToBool(vpc.IsDefault) {
			name += "*"
		}

		vpcCidr := aws.ToString(vpc.CidrBlock)
		v4Cidrs := []string{vpcCidr}

		for _, cbas := range vpc.CidrBlockAssociationSet {
			cbasCidr := aws.ToString(cbas.CidrBlock)
			if cbasCidr != vpcCidr && string(cbas.CidrBlockState.State) == "associated" {
				v4Cidrs = append(v4Cidrs, cbasCidr)
			}
		}

		v6Cidrs := make([]string, 0, len(vpc.Ipv6CidrBlockAssociationSet))
		for _, cba := range vpc.Ipv6CidrBlockAssociationSet {
			v6Cidrs = append(v6Cidrs, aws.ToString(cba.Ipv6CidrBlock))
		}

		if len(v6Cidrs) == 0 {
			v6Cidrs = []string{"-"}
		} else if len(v4Cidrs) > 1 || len(v6Cidrs) > 1 {
			p.multipleCidrs = true
		}

		table.AppendRow([]string{
			aws.ToString(vpc.VpcId),
			name,
			strings.Join(v4Cidrs, "\n"),
			strings.Join(v6Cidrs, "\n"),
		})
	}

	if p.multipleCidrs {
		table.SeparateRows()
	}

	table.Print(writer)
	if len(p.vpcs) > 0 {
		fmt.Println("* Indicates default VPC")
	}

	return nil
}

func (p *VpcPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.vpcs)
}

func (p *VpcPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.vpcs)
}

func (p *VpcPrinter) getVpcName(vpc types.Vpc) string {
	for _, tag := range vpc.Tags {
		if aws.ToString(tag.Key) == "Name" {
			return aws.ToString(tag.Value)
		}
	}
	return ""
}
