package ec2_instance

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

const maxNameLength = 30

type EC2Printer struct {
	reservations []types.Reservation
}

func NewPrinter(reservations []types.Reservation) *EC2Printer {
	return &EC2Printer{reservations}
}

func (p *EC2Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name", "Type", "Zone"})
	spot := 0

	for _, res := range p.reservations {
		for _, i := range res.Instances {
			age := durafmt.ParseShort(time.Since(aws.ToTime(i.LaunchTime)))

			instanceType := string(i.InstanceType)
			if string(i.InstanceLifecycle) == "spot" {
				instanceType = "*" + instanceType
				spot++
			}

			table.AppendRow([]string{
				age.String(),
				string(i.State.Name),
				aws.ToString(i.InstanceId),
				p.getInstanceName(i),
				instanceType,
				aws.ToString(i.Placement.AvailabilityZone),
			})
		}
	}

	table.Print(writer)
	if spot > 0 {
		fmt.Println("* Indicates Spot Instance")
	}

	return nil
}

func (p *EC2Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.reservations)
}

func (p *EC2Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.reservations)
}

func (p *EC2Printer) getInstanceName(instance types.Instance) string {
	name := ""
	for _, tag := range instance.Tags {
		if aws.ToString(tag.Key) == "Name" {
			name = aws.ToString(tag.Value)

			if len(name) > maxNameLength {
				name = name[:maxNameLength-3] + "..."
			}
			continue
		}
	}
	return name
}
