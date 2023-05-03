package cluster

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type ClusterPrinter struct {
	clusters   []*types.Cluster
	clusterURL string
}

func NewPrinter(clusters []*types.Cluster, clusterURL string) *ClusterPrinter {
	return &ClusterPrinter{clusters, clusterURL}
}

func (p *ClusterPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Cluster", "Version", "Platform", "Endpoint"})
	currentContext := false

	for _, cluster := range p.clusters {
		var endpoint string

		vpcConf := cluster.ResourcesVpcConfig
		if vpcConf == nil {
			endpoint = "-"
		} else if vpcConf.EndpointPublicAccess && !vpcConf.EndpointPrivateAccess {
			endpoint = "Public"
		} else if vpcConf.EndpointPublicAccess && vpcConf.EndpointPrivateAccess {
			endpoint = "Public/Private"
		} else {
			endpoint = "Private"
		}

		age := durafmt.ParseShort(time.Since(aws.ToTime(cluster.CreatedAt)))
		name := aws.ToString(cluster.Name)

		if aws.ToString(cluster.Endpoint) == p.clusterURL {
			currentContext = true
			name = "*" + name
		}

		table.AppendRow([]string{
			age.String(),
			string(cluster.Status),
			name,
			aws.ToString(cluster.Version),
			aws.ToString(cluster.PlatformVersion),
			endpoint,
		})
	}

	table.Print(writer)
	if currentContext {
		fmt.Println("* Indicates current context in local kubeconfig")
	}

	return nil
}

func (p *ClusterPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.clusters)
}

func (p *ClusterPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.clusters)
}
