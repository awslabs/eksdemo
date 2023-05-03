package nodegroup

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type NodegroupPrinter struct {
	Nodegroups []*types.Nodegroup
}

func NewPrinter(Nodegroups []*types.Nodegroup) *NodegroupPrinter {
	return &NodegroupPrinter{Nodegroups}
}

func (p *NodegroupPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Nodes", "Min", "Max", "Version", "Type", "Instance(s)"})

	for _, n := range p.Nodegroups {
		age := durafmt.ParseShort(time.Since(aws.ToTime(n.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(n.Status),
			aws.ToString(n.NodegroupName),
			strconv.Itoa(int(aws.ToInt32(n.ScalingConfig.DesiredSize))),
			strconv.Itoa(int(aws.ToInt32(n.ScalingConfig.MinSize))),
			strconv.Itoa(int(aws.ToInt32(n.ScalingConfig.MaxSize))),
			aws.ToString(n.ReleaseVersion),
			string(n.CapacityType),
			strings.Join(n.InstanceTypes, ","),
		})
	}

	table.Print(writer)

	return nil
}

func (p *NodegroupPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Nodegroups)
}

func (p *NodegroupPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Nodegroups)
}
