package network_acl_rule

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type NetworkAclRulePrinter struct {
	naclRules []types.NetworkAclEntry
}

func NewPrinter(naclRules []types.NetworkAclEntry) *NetworkAclRulePrinter {
	return &NetworkAclRulePrinter{naclRules}
}

func (p *NetworkAclRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Rule #", "Protocol", "Port Range", "Source", "Action"})
	table.SetColumnAlignment([]int{
		printer.ALIGN_RIGHT, printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_LEFT,
	})

	for _, r := range p.naclRules {
		ruleNumber := strconv.Itoa(int(aws.ToInt32(r.RuleNumber)))
		if ruleNumber == "32767" || ruleNumber == "32768" {
			ruleNumber = "*"
		}

		source := ""
		if cidr := aws.ToString(r.CidrBlock); cidr != "" {
			source = cidr
		}
		if v6cidr := aws.ToString(r.Ipv6CidrBlock); v6cidr != "" {
			source = v6cidr
		}

		table.AppendRow([]string{
			ruleNumber,
			toProtocolString(r),
			toPortRangeString(r.PortRange, r.IcmpTypeCode),
			source,
			strings.Title(string(r.RuleAction)),
		})
	}

	table.Print(writer)
	if len(p.naclRules) > 0 {
		if aws.ToBool(p.naclRules[0].Egress) {
			fmt.Println("* Displaying egress rules")
		} else {
			fmt.Println("* Displaying ingress rules, use --egress for egress rules")
		}
	}

	return nil
}

func (p *NetworkAclRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.naclRules)
}

func (p *NetworkAclRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.naclRules)
}

func toProtocolString(rule types.NetworkAclEntry) string {
	proto := aws.ToString(rule.Protocol)
	switch {
	case proto == "-1":
		return "All"
	case proto == "1":
		return "ICMP (1)"
	case proto == "6":
		return "TCP (6)"
	case proto == "17":
		return "UDP (17)"
	case proto == "58":
		return "IPv6-ICMP (58)"
	default:
		return proto
	}
}

func toPortRangeString(pr *types.PortRange, itc *types.IcmpTypeCode) string {
	switch {
	case pr == nil && itc == nil:
		return "All"
	case itc != nil && aws.ToInt32(itc.Code) == -1:
		return "All"
	case itc != nil:
		return fmt.Sprintf("Code: %d, Type: %d", aws.ToInt32(itc.Code), aws.ToInt32(itc.Type))
	case aws.ToInt32(pr.From) == aws.ToInt32(pr.To):
		return strconv.Itoa(int(aws.ToInt32(pr.From)))
	default:
		return strconv.Itoa(int(aws.ToInt32(pr.From))) + "-" + strconv.Itoa(int(aws.ToInt32(pr.To)))
	}
}
