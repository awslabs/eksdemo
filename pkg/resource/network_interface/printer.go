package network_interface

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type NetworkInterfacePrinter struct {
	accountId         string
	networkInterfaces []types.NetworkInterface
}

func NewPrinter(networkInterfaces []types.NetworkInterface, accountId string) *NetworkInterfacePrinter {
	return &NetworkInterfacePrinter{accountId, networkInterfaces}
}

func (p *NetworkInterfacePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Instance Id or...", "Private IPv4", "IPs", "SGs", "Subnet"})

	for _, eni := range p.networkInterfaces {
		id := aws.ToString(eni.NetworkInterfaceId)
		instanceId := ""

		if eni.Attachment == nil {
			instanceId = "detached"
		} else {
			if aws.ToInt32(eni.Attachment.DeviceIndex) == 0 {
				id = "*" + id
			}

			if aws.ToString(eni.Attachment.InstanceId) != "" {
				instanceId = aws.ToString(eni.Attachment.InstanceId)
			} else if string(eni.InterfaceType) != "interface" {
				instanceId = string(eni.InterfaceType)
			} else if aws.ToString(eni.Attachment.InstanceOwnerId) == "amazon-elb" {
				instanceId = "load-balancer"

				// Identify EKS Control Plane cross-account and Fargate ENIs
				// This identifies any interface ENI that is attached to EC2 but owned by a different account
			} else if aws.ToString(eni.Attachment.InstanceOwnerId) != p.accountId {
				instanceId = "eks_control_plane"
				if !strings.HasPrefix(aws.ToString(eni.Description), "Amazon EKS") {
					instanceId = "fargate_pod"
				}
			}
		}

		table.AppendRow([]string{
			id,
			instanceId,
			aws.ToString(eni.PrivateIpAddress),
			strconv.Itoa(len(eni.PrivateIpAddresses)),
			strconv.Itoa(len(eni.Groups)),
			aws.ToString(eni.SubnetId),
		})
	}

	table.Print(writer)
	if len(p.networkInterfaces) > 0 {
		fmt.Println("* Indicates Primary network interface")
	}

	return nil
}

func (p *NetworkInterfacePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.networkInterfaces)
}

func (p *NetworkInterfacePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.networkInterfaces)
}
