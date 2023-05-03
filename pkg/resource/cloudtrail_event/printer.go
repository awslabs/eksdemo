package cloudtrail_event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

const maxNameLength = 30

type CloudtrailEventPrinter struct {
	events []types.Event
}

func NewPrinter(events []types.Event) *CloudtrailEventPrinter {
	return &CloudtrailEventPrinter{events}
}

func (p *CloudtrailEventPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()

	table.SetHeader([]string{"Age", "Id", "Name", "Username"})
	table.SetColumnAlignment([]int{
		printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_LEFT,
	})

	for _, e := range p.events {
		age := durafmt.ParseShort(time.Since(aws.ToTime(e.EventTime)))

		name := aws.ToString(e.EventName) + " (" + strings.Split(aws.ToString(e.EventSource), ".")[0] + ")"
		username := aws.ToString(e.Username)

		if len(username) > maxNameLength {
			username = username[:maxNameLength-3] + "..."
		}

		table.AppendRow([]string{
			age.String(),
			aws.ToString(e.EventId),
			name,
			username,
		})

	}
	table.Print(writer)

	return nil
}

func (p *CloudtrailEventPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.events)
}

func (p *CloudtrailEventPrinter) PrintYAML(writer io.Writer) error {
	if err := p.prettyPrintPolicyDocuments(); err != nil {
		return err
	}

	return printer.EncodeYAML(writer, p.events)
}

func (p *CloudtrailEventPrinter) prettyPrintPolicyDocuments() error {
	var prettyJSON bytes.Buffer

	for i, event := range p.events {
		err := json.Indent(&prettyJSON, []byte(aws.ToString(event.CloudTrailEvent)), "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON for event id %q: %w", aws.ToString(event.EventId), err)
		}

		p.events[i].CloudTrailEvent = aws.String(prettyJSON.String())
	}

	return nil
}
