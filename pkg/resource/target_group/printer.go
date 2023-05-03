package target_group

import (
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type TargetGroupPrinter struct {
	targetGroups []types.TargetGroup
}

func NewPrinter(targetGroups []types.TargetGroup) *TargetGroupPrinter {
	return &TargetGroupPrinter{targetGroups}
}

func (p *TargetGroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Proto:Port", "Load Balancer"})

	for _, tg := range p.targetGroups {
		proto := string(tg.Protocol)
		if version := aws.ToString(tg.ProtocolVersion); version != "" {
			proto = version
		}

		table.AppendRow([]string{
			aws.ToString(tg.TargetGroupName),
			string(tg.TargetType),
			proto + ":" + strconv.Itoa(int(aws.ToInt32(tg.Port))),
			getLoadBalancer(tg),
		})
	}

	table.Print(writer)

	return nil
}

func (p *TargetGroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.targetGroups)
}

func (p *TargetGroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.targetGroups)
}

func getLoadBalancer(tg types.TargetGroup) string {
	lbArns := tg.LoadBalancerArns
	if len(lbArns) == 0 {
		return "None associated"
	} else if len(lbArns) > 1 {
		return strconv.Itoa(len(lbArns)) + " LBs associated"
	}

	splitArn := strings.Split(lbArns[0], "/")
	if len(splitArn) < 2 {
		return "Error parsing LoadBalancer ARN"
	}
	return splitArn[2]
}
