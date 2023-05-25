package load_balancer

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/awslabs/eksdemo/pkg/printer"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hako/durafmt"
)

type LoadBalancerPrinter struct {
	*LoadBalancers
}

func NewPrinter(loadBalancers *LoadBalancers) *LoadBalancerPrinter {
	return &LoadBalancerPrinter{loadBalancers}
}

func (p *LoadBalancerPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Name", "Type", "Stack", "AZs", "SGs"})

	for _, lb := range p.V1 {
		age := durafmt.ParseShort(time.Since(aws.ToTime(lb.CreatedTime)))

		name := aws.ToString(lb.LoadBalancerName)
		if aws.ToString(lb.Scheme) == "internal" {
			name = "*" + name
		}

		table.AppendRow([]string{
			age.String(),
			"-",
			name,
			"CLB",
			"ipv4",
			strconv.Itoa(len(lb.AvailabilityZones)),
			strconv.Itoa(len(lb.SecurityGroups)),
		})
	}

	for _, lb := range p.V2 {
		age := durafmt.ParseShort(time.Since(aws.ToTime(lb.CreatedTime)))

		name := aws.ToString(lb.LoadBalancerName)
		if string(lb.Scheme) == "internal" {
			name = "*" + name
		}

		elbType := "unknown"
		switch string(lb.Type) {
		case "application":
			elbType = "ALB"
		case "network":
			elbType = "NLB"
		case "gateway":
			elbType = "GWLB"
		}

		table.AppendRow([]string{
			age.String(),
			string(lb.State.Code),
			name,
			elbType,
			string(lb.IpAddressType),
			strconv.Itoa(len(lb.AvailabilityZones)),
			strconv.Itoa(len(lb.SecurityGroups)),
		})
	}

	table.Print(writer)
	if len(p.V1) > 0 || len(p.V2) > 0 {
		fmt.Println("* Indicates internal load balancer")
	}

	return nil
}

func (p *LoadBalancerPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.LoadBalancers)
}

func (p *LoadBalancerPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.LoadBalancers)
}
