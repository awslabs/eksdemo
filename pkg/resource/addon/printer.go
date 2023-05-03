package addon

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type AddonPrinter struct {
	addons []*types.Addon
}

func NewPrinter(addons []*types.Addon) *AddonPrinter {
	return &AddonPrinter{addons}
}

func (p *AddonPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Version"})

	for _, addon := range p.addons {
		age := durafmt.ParseShort(time.Since(aws.ToTime(addon.CreatedAt)))
		name := aws.ToString(addon.AddonName)

		table.AppendRow([]string{
			age.String(),
			string(addon.Status),
			name,
			aws.ToString(addon.AddonVersion),
		})
	}

	table.Print(writer)

	return nil
}

func (p *AddonPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.addons)
}

func (p *AddonPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.addons)
}
