package cloudtrail_trail

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type TrailPrinter struct {
	trails []Trail
}

func NewPrinter(trails []Trail) *TrailPrinter {
	return &TrailPrinter{trails}
}

func (p *TrailPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"\nAge", "\nName", "Home\nRegion", "Multi-\nRegion", "Org\nTrail", "\nStatus"})

	for _, p := range p.trails {
		age := durafmt.ParseShort(time.Since(aws.ToTime(p.TrailStatus.StartLoggingTime)))

		multiRegion := "No"
		if aws.ToBool(p.Trail.IsMultiRegionTrail) {
			multiRegion = "Yes"
		}

		orgTrail := "No"
		if aws.ToBool(p.Trail.IsOrganizationTrail) {
			orgTrail = "Yes"
		}

		status := "Logging"
		if !aws.ToBool(p.TrailStatus.IsLogging) {
			status = "Off"
		}

		table.AppendRow([]string{
			age.String(),
			aws.ToString(p.Trail.Name),
			aws.ToString(p.Trail.HomeRegion),
			multiRegion,
			orgTrail,
			status,
		})

	}
	table.Print(writer)

	return nil
}

func (p *TrailPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.trails)
}

func (p *TrailPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.trails)
}
