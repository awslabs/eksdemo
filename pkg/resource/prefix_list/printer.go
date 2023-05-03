package prefix_list

import (
	"io"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type PrefixListPrinter struct {
	prefixLists []PrefixList
}

func NewPrinter(prefixLists []PrefixList) *PrefixListPrinter {
	return &PrefixListPrinter{prefixLists}
}

func (p *PrefixListPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()

	table.SetHeader([]string{"Id", "Name", "Max", "Owner Id"})
	table.SetColumnAlignment([]int{
		printer.ALIGN_LEFT, printer.ALIGN_LEFT, printer.ALIGN_RIGHT, printer.ALIGN_RIGHT,
	})

	for _, pl := range p.prefixLists {
		maxEntries := "-"
		if pl.PrefixList.MaxEntries != nil {
			maxEntries = strconv.Itoa(int(aws.ToInt32(pl.PrefixList.MaxEntries)))
		}

		table.AppendRow([]string{
			aws.ToString(pl.PrefixList.PrefixListId),
			aws.ToString(pl.PrefixList.PrefixListName),
			maxEntries,
			aws.ToString(pl.PrefixList.OwnerId),
		})
	}

	table.Print(writer)

	return nil
}

func (p *PrefixListPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.prefixLists)
}

func (p *PrefixListPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.prefixLists)
}
