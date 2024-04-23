package session

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type SessionPrinter struct {
	sessions []types.Session
}

func NewPrinter(sessions []types.Session) *SessionPrinter {
	return &SessionPrinter{sessions}
}

func (p *SessionPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Instance"})

	for _, s := range p.sessions {
		age := durafmt.ParseShort(time.Since(aws.ToTime(s.StartDate)))

		table.AppendRow([]string{
			age.String(),
			string(s.Status),
			aws.ToString(s.SessionId),
			aws.ToString(s.Target),
		})
	}

	table.Print(writer)

	return nil
}

func (p *SessionPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.sessions)
}

func (p *SessionPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.sessions)
}
