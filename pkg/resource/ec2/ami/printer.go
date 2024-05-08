package ami

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	amis []types.Image
}

func NewPrinter(amis []types.Image) *Printer {
	return &Printer{amis}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Id", "Name"})

	for _, ami := range p.amis {
		creation, err := time.Parse(time.RFC3339, aws.ToString(ami.CreationDate))

		if err != nil {
			return fmt.Errorf("failed to parse CreationDate: %w", err)
		}

		age := durafmt.ParseShort(time.Since(creation))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(ami.ImageId),
			aws.ToString(ami.Name),
		})
	}

	table.Print(writer)

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.amis)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.amis)
}
