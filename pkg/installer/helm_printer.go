package installer

import (
	"fmt"
	"io"
	"strconv"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type HelmPrinter struct {
	release *HelmInstaller
	options application.Options
	values  string
}

func NewHelmPrinter(release *HelmInstaller, options application.Options, values string) *HelmPrinter {
	return &HelmPrinter{release, options, values}
}

func (p *HelmPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()

	rel := p.release
	opt := p.options.Common()

	table.AppendRow([]string{"Application Version", opt.Version})
	table.AppendRow([]string{"Chart Version", opt.ChartVersion})
	table.AppendRow([]string{"Chart Repository", rel.RepositoryURL})
	table.AppendRow([]string{"Chart Name", rel.ChartName})
	table.AppendRow([]string{"Release Name", rel.ReleaseName})
	table.AppendRow([]string{"Namespace", opt.Namespace})
	table.AppendRow([]string{"Wait", strconv.FormatBool(rel.Wait)})

	table.Print(writer)

	fmt.Printf("Set Values: %v\n", opt.SetValues)
	fmt.Println("Values File:")
	fmt.Println(p.values)

	return nil
}
