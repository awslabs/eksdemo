package event_rule

import (
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type EventBridgeRulePrinter struct {
	rules []types.Rule
}

func NewPrinter(rules []types.Rule) *EventBridgeRulePrinter {
	return &EventBridgeRulePrinter{rules}
}

func (p *EventBridgeRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Status", "Name", "Type"})

	for _, r := range p.rules {
		eventType := "Standard"
		if r.ManagedBy != nil {
			eventType = "Managed"
		}

		table.AppendRow([]string{
			string(r.State),
			aws.ToString(r.Name),
			eventType,
		})
	}

	table.Print(writer)

	return nil
}

func (p *EventBridgeRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.rules)
}

func (p *EventBridgeRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.rules)
}
