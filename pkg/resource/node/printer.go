package node

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
	v1 "k8s.io/api/core/v1"
)

type NodePrinter struct {
	nodes []v1.Node
}

func NewPrinter(nodes []v1.Node) *NodePrinter {
	return &NodePrinter{nodes}
}

func (p *NodePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Instance Id", "Type", "Zone", "Nodegroup"})

	for _, node := range p.nodes {
		age := durafmt.ParseShort(time.Since(node.CreationTimestamp.Time))
		name := strings.Split(node.Name, ".")[0] + ".*"

		instanceId := node.Spec.ProviderID[strings.LastIndex(node.Spec.ProviderID, "/")+1:]
		if !strings.HasPrefix(instanceId, "i-") {
			instanceId = "-"
		}

		labels := node.GetLabels()

		instanceType, ok := labels["node.kubernetes.io/instance-type"]
		if !ok {
			instanceType, ok = labels["eks.amazonaws.com/compute-type"]
			if !ok {
				instanceType = "unknown"
			}
		}

		nodegroup, ok := labels["eks.amazonaws.com/nodegroup"]
		if !ok {
			nodegroup = "-"
		}

		zone, ok := labels["topology.kubernetes.io/zone"]
		if !ok {
			zone = "unknown"
		}

		table.AppendRow([]string{
			age.String(),
			name,
			instanceId,
			instanceType,
			zone,
			nodegroup,
		})
	}

	table.Print(writer)
	if len(p.nodes) > 0 {
		node := p.nodes[0]
		nodeSuffix := node.Name[strings.Index(node.Name, ".")+1:]
		fmt.Printf("* Names end with %q\n", nodeSuffix)
	}

	return nil
}

func (p *NodePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.nodes)
}

func (p *NodePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.nodes)
}
