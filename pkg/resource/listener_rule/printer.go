package listener_rule

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource/listener"
)

type ListenerRulePrinter struct {
	rules []types.Rule
}

func NewPrinter(rules []types.Rule) *ListenerRulePrinter {
	return &ListenerRulePrinter{rules}
}

func (p *ListenerRulePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Prior", "Conditions", "Actions"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

	for _, r := range p.rules {
		table.AppendRow([]string{
			resourceId.FindString(aws.ToString(r.RuleArn)),
			aws.ToString(r.Priority),
			strings.Join(printConditions(r.Conditions), "\n"),
			strings.Join(listener.PrintActions(r.Actions), "\n"),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *ListenerRulePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.rules)
}

func (p *ListenerRulePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.rules)
}

func printConditions(elbConditions []types.RuleCondition) (conditions []string) {
	if len(elbConditions) == 0 {
		conditions = []string{"Requests otherwise not routed"}
	} else {
		for _, c := range elbConditions {
			switch {
			case c.HostHeaderConfig != nil:
				hostHeaders := c.HostHeaderConfig.Values
				conditions = append(conditions, "Host is "+strings.Join(hostHeaders, " OR "))

			case c.HttpHeaderConfig != nil:
				header := aws.ToString(c.HttpHeaderConfig.HttpHeaderName)
				headerValues := c.HttpHeaderConfig.Values
				conditions = append(conditions, "Http header "+header+" is "+strings.Join(headerValues, " OR "))

			case c.HttpRequestMethodConfig != nil:
				methods := c.HttpRequestMethodConfig.Values
				conditions = append(conditions, "Http request method is "+strings.Join(methods, " OR "))

			case c.PathPatternConfig != nil:
				pathPattern := c.PathPatternConfig.Values
				conditions = append(conditions, "Path is "+strings.Join(pathPattern, " OR "))

			case c.QueryStringConfig != nil:
				kvText := []string{}
				for _, kvp := range c.QueryStringConfig.Values {
					condText := fmt.Sprintf("%s:%s", aws.ToString(kvp.Key), aws.ToString(kvp.Value))
					kvText = append(kvText, condText)
				}
				conditions = append(conditions, "Query string is "+strings.Join(kvText, " OR "))

			case c.SourceIpConfig != nil:
				sourceIPs := c.SourceIpConfig.Values
				conditions = append(conditions, "Source IP is "+strings.Join(sourceIPs, " OR "))
			}
		}
	}
	return
}
