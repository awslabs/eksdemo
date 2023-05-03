package volume

import (
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

const maxNameLength = 25

type VolumePrinter struct {
	volumes []types.Volume
}

func NewPrinter(volumes []types.Volume) *VolumePrinter {
	return &VolumePrinter{volumes}
}

func (p *VolumePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "State", "Id", "Name", "Type", "GiB", "AZ"})

	for _, v := range p.volumes {
		age := durafmt.ParseShort(time.Since(aws.ToTime(v.CreateTime)))

		table.AppendRow([]string{
			age.String(),
			string(v.State),
			aws.ToString(v.VolumeId),
			p.getVolumeName(v),
			string(v.VolumeType),
			strconv.Itoa(int(aws.ToInt32(v.Size))),
			aws.ToString(v.AvailabilityZone),
		})
	}
	table.Print(writer)

	return nil
}

func (p *VolumePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.volumes)
}

func (p *VolumePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.volumes)
}

func (p *VolumePrinter) getVolumeName(volume types.Volume) string {
	name := ""
	for _, tag := range volume.Tags {
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
