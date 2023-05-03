package amp_rule

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type AmpRulePrinter struct {
	ruleGroupNamespaces []*types.RuleGroupsNamespaceDescription
}

func NewPrinter(ruleGroupNamespaces []*types.RuleGroupsNamespaceDescription) *AmpRulePrinter {
	return &AmpRulePrinter{ruleGroupNamespaces}
}

func (p *AmpRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Namespace"})

	for _, rgn := range p.ruleGroupNamespaces {
		age := durafmt.ParseShort(time.Since(aws.ToTime(rgn.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(rgn.Status.StatusCode),
			aws.ToString(rgn.Name),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AmpRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.ruleGroupNamespaces)
}

func (p *AmpRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.ruleGroupNamespaces)
}
