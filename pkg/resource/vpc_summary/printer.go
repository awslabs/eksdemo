package vpc_summary

import (
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type VpcSummaryPrinter struct {
	*VpcSummary
	vpcFilter bool
	showIds   bool
}

func NewPrinter(vpcSummary *VpcSummary, vpcFilter, showIds bool) *VpcSummaryPrinter {
	return &VpcSummaryPrinter{vpcSummary, vpcFilter, showIds}
}

func (p *VpcSummaryPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	header := []string{"Resource", "Count"}

	sum := p.VpcSummary
	eip := []string{"Elastic IPs", strconv.Itoa(len(sum.eips))}
	endpoint := []string{"Endpoints", strconv.Itoa(len(sum.endpoints))}
	ec2 := []string{"Instances (Running)", strconv.Itoa(len(sum.instances))}
	igw := []string{"Internet Gateways", strconv.Itoa(len(sum.internetGWs))}
	lb := []string{"Load Balancers", strconv.Itoa(len(sum.loadBalancers.V1) + len(sum.loadBalancers.V2))}
	nat := []string{"NAT Gateways", strconv.Itoa(len(sum.natGWs))}
	rt := []string{"Route Tables", strconv.Itoa(len(sum.routeTables))}
	sg := []string{"Security Groups", strconv.Itoa(len(sum.securityGroups))}
	subnet := []string{"Subnets", strconv.Itoa(len(sum.subnets))}
	vpc := []string{"VPCs", strconv.Itoa(len(sum.vpcs))}

	if p.showIds {
		header = append(header, "Id(s)")
		table.SeparateRows()

		ec2Ids := make([]string, 0, len(sum.instances))
		for _, r := range sum.instances {
			for _, i := range r.Instances {
				ec2Ids = append(ec2Ids, aws.ToString(i.InstanceId))
			}
		}

		elasticIpIds := make([]string, 0, len(sum.eips))
		for _, e := range sum.eips {
			elasticIpIds = append(elasticIpIds, aws.ToString(e.AllocationId))
		}

		endpointIds := make([]string, 0, len(sum.endpoints))
		for _, e := range sum.endpoints {
			endpointIds = append(endpointIds, aws.ToString(e.VpcEndpointId))
		}

		igwIds := make([]string, 0, len(sum.internetGWs))
		for _, i := range sum.internetGWs {
			igwIds = append(igwIds, aws.ToString(i.InternetGatewayId))
		}

		lbIds := make([]string, 0, len(sum.loadBalancers.V1)+len(sum.loadBalancers.V2))
		for _, v1 := range sum.loadBalancers.V1 {
			lbIds = append(lbIds, aws.ToString(v1.LoadBalancerName))
		}
		for _, v2 := range sum.loadBalancers.V2 {
			lbIds = append(lbIds, aws.ToString(v2.LoadBalancerName))
		}

		natIds := make([]string, 0, len(sum.natGWs))
		for _, n := range sum.natGWs {
			natIds = append(natIds, aws.ToString(n.NatGatewayId))
		}

		rtIds := make([]string, 0, len(sum.routeTables))
		for _, rt := range sum.routeTables {
			rtIds = append(rtIds, aws.ToString(rt.RouteTableId))
		}

		sgIds := make([]string, 0, len(sum.securityGroups))
		for _, s := range sum.securityGroups {
			sgIds = append(sgIds, aws.ToString(s.GroupId))
		}

		subnetIds := make([]string, 0, len(sum.subnets))
		for _, s := range sum.subnets {
			subnetIds = append(subnetIds, aws.ToString(s.SubnetId))
		}

		vpcIds := make([]string, 0, len(sum.vpcs))
		for _, vpc := range sum.vpcs {
			vpcIds = append(vpcIds, aws.ToString(vpc.VpcId))
		}

		ec2 = append(ec2, strings.Join(ec2Ids, ", "))
		eip = append(eip, strings.Join(elasticIpIds, ", "))
		endpoint = append(endpoint, strings.Join(endpointIds, ", "))
		igw = append(igw, strings.Join(igwIds, ", "))
		lb = append(lb, strings.Join(lbIds, ", "))
		nat = append(nat, strings.Join(natIds, ", "))
		rt = append(rt, strings.Join(rtIds, ", "))
		sg = append(sg, strings.Join(sgIds, ", "))
		subnet = append(subnet, strings.Join(subnetIds, ", "))
		vpc = append(vpc, strings.Join(vpcIds, ", "))
	}

	table.SetHeader(header)
	if !p.vpcFilter {
		table.AppendRow(eip)

	}
	table.AppendRow(endpoint)
	table.AppendRow(ec2)
	table.AppendRow(igw)
	table.AppendRow(lb)
	table.AppendRow(nat)
	table.AppendRow(rt)
	table.AppendRow(sg)
	table.AppendRow(subnet)
	table.AppendRow(vpc)
	table.Print(writer)

	return nil
}

func (p *VpcSummaryPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.vpcs)
}

func (p *VpcSummaryPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.vpcs)
}

// func (p *VpcSummaryPrinter) getVpcName(vpc types.Vpc) string {
// 	for _, tag := range vpc.Tags {
// 		if aws.ToString(tag.Key) == "Name" {
// 			return aws.ToString(tag.Value)
// 		}
// 	}
// 	return ""
// }
