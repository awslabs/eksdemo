package cloudformation_stack

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

const maxNameLength = 70

type CloudFormationPrinter struct {
	Stacks []types.Stack
}

func NewPrinter(stacks []types.Stack) *CloudFormationPrinter {
	return &CloudFormationPrinter{stacks}
}

func (p *CloudFormationPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name"})

	for _, s := range p.Stacks {

		age := durafmt.ParseShort(time.Since(aws.ToTime(s.CreationTime)))
		name := aws.ToString(s.StackName)

		if len(name) > maxNameLength {
			name = name[:maxNameLength-3] + "..."
		}

		table.AppendRow([]string{
			age.String(),
			string(s.StackStatus),
			name,
		})
	}

	table.Print(writer)

	return nil
}

func (p *CloudFormationPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.Stacks)
}

func (p *CloudFormationPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.Stacks)
}
