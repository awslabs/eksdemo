package amp_workspace

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type AmpWorkspacePrinter struct {
	Workspaces []AmpWorkspace
}

func NewPrinter(workspaces []AmpWorkspace) *AmpWorkspacePrinter {
	return &AmpWorkspacePrinter{workspaces}
}

func (p *AmpWorkspacePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Alias", "Workspace Id"})

	for _, w := range p.Workspaces {

		age := durafmt.ParseShort(time.Since(aws.ToTime(w.Workspace.CreatedAt)))

		table.AppendRow([]string{
			age.String(),
			string(w.Workspace.Status.StatusCode),
			aws.ToString(w.Workspace.Alias),
			aws.ToString(w.Workspace.WorkspaceId),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AmpWorkspacePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Workspaces)
}

func (p *AmpWorkspacePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Workspaces)
}
