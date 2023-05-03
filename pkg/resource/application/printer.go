package application

import (
	"fmt"
	"io"

	"github.com/awslabs/eksdemo/pkg/printer"
	"helm.sh/helm/v3/pkg/release"
)

type ApplicationPrinter struct {
	releases []*release.Release
}

func NewPrinter(releases []*release.Release) *ApplicationPrinter {
	return &ApplicationPrinter{releases}
}

func (p *ApplicationPrinter) PrintTable(writer io.Writer) error {
	if len(p.releases) == 0 {
		fmt.Fprint(writer, "No helm releases found.\n")
		return nil
	}

	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Namespace", "Version", "Status", "Chart"})

	for _, r := range p.releases {
		table.AppendRow([]string{
			r.Name,
			r.Namespace,
			r.Chart.Metadata.AppVersion,
			r.Info.Status.String(),
			r.Chart.Metadata.Version,
		})
	}
	table.Print(writer)

	return nil
}

func (p *ApplicationPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.releases)
}

func (p *ApplicationPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.releases)
}
