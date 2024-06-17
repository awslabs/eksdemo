package printer

import (
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TablePrinter struct {
	alignment    []int
	header       []string
	data         [][]string
	noTextWrap   bool
	rowSeperator bool
}

const (
	ALIGN_DEFAULT = iota
	ALIGN_CENTER
	ALIGN_RIGHT
	ALIGN_LEFT
)

// TODO: use make to set the size of the slice
func NewTablePrinter() *TablePrinter {
	return &TablePrinter{}
}

func (p *TablePrinter) AppendRow(row []string) {
	p.data = append(p.data, row)
}

func (p *TablePrinter) NoTextWrap() {
	p.noTextWrap = true
}

func (p *TablePrinter) SetColumnAlignment(keys []int) {
	p.alignment = keys
}

func (p *TablePrinter) SetHeader(header []string) {
	p.header = header
}

func (p *TablePrinter) SeparateRows() {
	p.rowSeperator = true
}

func (p *TablePrinter) Print(writer io.Writer) {
	if len(p.data) == 0 {
		fmt.Println("No resources found.")
		return
	}
	table := tablewriter.NewWriter(writer)
	table.SetAutoFormatHeaders(false)

	if len(p.alignment) > 0 {
		table.SetColumnAlignment(p.alignment)
	}

	if p.noTextWrap {
		table.SetAutoWrapText(false)
	}

	if p.rowSeperator {
		table.SetRowLine(true)
	}

	table.SetHeader(p.header)
	table.AppendBulk(p.data)
	table.Render()
}

func (p *TablePrinter) TruncateBeginning(text string, max int) string {
	if len(text) > max {
		return "..." + text[len(text)-max+3:]
	}
	return text
}

func (p *TablePrinter) TruncateMiddle(text string, max int) string {
	if len(text) > max {
		return text[:max/2-1] + "..." + text[len(text)-max/2+1:]
	}
	return text
}

func (p *TablePrinter) TruncateMiddleWithEllipsisLocation(text string, ellipisLocation, max int) string {
	if len(text) > max {
		return text[:ellipisLocation-1] + "..." + text[len(text)-max+ellipisLocation+2:]
	}
	return text
}
