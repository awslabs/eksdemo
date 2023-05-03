package security_group_rule

import (
	"fmt"
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type SecurityGroupRulePrinter struct {
	securityGroupRules []types.SecurityGroupRule
	clusterName        string
}

func NewPrinter(securityGroupRules []types.SecurityGroupRule, clusterName string) *SecurityGroupRulePrinter {
	return &SecurityGroupRulePrinter{securityGroupRules, clusterName}
}

func (p *SecurityGroupRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Rule Id", "Proto", "Ports", "Source", "Description"})

	for _, sgr := range p.securityGroupRules {
		fromPort := aws.ToInt32(sgr.FromPort)
		toPort := aws.ToInt32(sgr.ToPort)

		id := aws.ToString(sgr.SecurityGroupRuleId)
		if aws.ToBool(sgr.IsEgress) {
			id = "*" + id
		}

		ports := "All"
		if fromPort != -1 {
			if fromPort == toPort {
				ports = strconv.Itoa(int(fromPort))
			} else {
				ports = strconv.Itoa(int(fromPort)) + "-" + strconv.Itoa(int(toPort))
			}
		}

		protocol := "All"
		if aws.ToString(sgr.IpProtocol) != "-1" {
			protocol = aws.ToString(sgr.IpProtocol)
		}

		source := "-"
		if sgr.ReferencedGroupInfo != nil {
			source = aws.ToString(sgr.ReferencedGroupInfo.GroupId)
		} else if sgr.CidrIpv4 != nil {
			source = aws.ToString(sgr.CidrIpv4)
		} else if sgr.PrefixListId != nil {
			source = aws.ToString(sgr.PrefixListId)
		}

		table.AppendRow([]string{
			id,
			protocol,
			ports,
			source,
			aws.ToString(sgr.Description),
		})
	}

	table.SeparateRows()
	table.Print(writer)
	if len(p.securityGroupRules) > 0 {
		fmt.Println("* Indicates egress rule")
	}

	return nil
}

func (p *SecurityGroupRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.securityGroupRules)
}

func (p *SecurityGroupRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.securityGroupRules)
}
