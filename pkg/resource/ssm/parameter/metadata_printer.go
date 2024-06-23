package parameter

import (
	"io"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type MetadataPrinter struct {
	params []types.ParameterMetadata
}

func NewMetadataPrinter(params []types.ParameterMetadata) *MetadataPrinter {
	return &MetadataPrinter{params}
}

func (p *MetadataPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Tier", "Type", "Ver"})

	for _, param := range p.params {
		age := durafmt.ParseShort(time.Since(aws.ToTime(param.LastModifiedDate)))

		table.AppendRow([]string{
			age.String(),
			table.TruncateMiddleWithEllipsisLocation(aws.ToString(param.Name), 20, 80),
			string(param.Tier),
			aws.ToString(param.DataType),
			strconv.Itoa(int(param.Version)),
		})
	}

	table.Print(writer)

	return nil
}

func (p *MetadataPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.params)
}

func (p *MetadataPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.params)
}
