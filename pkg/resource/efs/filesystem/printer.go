package filesystem

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
)

type Printer struct {
	fileSystems []types.FileSystemDescription
}

func NewPrinter(fileSystems []types.FileSystemDescription) *Printer {
	return &Printer{fileSystems}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Id", "Name", "Size"})

	for _, fs := range p.fileSystems {
		age := durafmt.ParseShort(time.Since(aws.ToTime(fs.CreationTime)))

		table.AppendRow([]string{
			age.String(),
			string(fs.LifeCycleState),
			aws.ToString(fs.FileSystemId),
			aws.ToString(fs.Name),
			humanize.IBytes(uint64(fs.SizeInBytes.Value)),
		})
	}
	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.fileSystems)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.fileSystems)
}
