package accessentry

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type Printer struct {
	accessEntries []*types.AccessEntry
}

func NewPrinter(accessEntries []*types.AccessEntry) *Printer {
	return &Printer{accessEntries}
}

func (p *Printer) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Principal ARN"})

	for _, ae := range p.accessEntries {
		age := durafmt.ParseShort(time.Since(aws.ToTime(ae.CreatedAt)))
		arn := aws.ToString(ae.PrincipalArn)

		table.AppendRow([]string{
			age.String(),
			"*" + arn[strings.LastIndex(arn, ":"):],
		})
	}

	table.Print(writer)
	if len(p.accessEntries) > 0 {
		arn := aws.ToString(p.accessEntries[0].PrincipalArn)
		prefix := arn[:strings.LastIndex(arn, ":")]
		fmt.Printf("* ARNs start with %q\n", prefix)
	}

	return nil
}

func (p *Printer) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.accessEntries)
}

func (p *Printer) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.accessEntries)
}
