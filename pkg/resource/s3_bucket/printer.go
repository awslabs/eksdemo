package s3_bucket

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type BucketPrinter struct {
	buckets []types.Bucket
}

func NewPrinter(buckets []types.Bucket) *BucketPrinter {
	return &BucketPrinter{buckets}
}

func (p *BucketPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name"})

	for _, b := range p.buckets {
		age := durafmt.ParseShort(time.Since(aws.ToTime(b.CreationDate)))

		table.AppendRow([]string{
			age.String(),
			aws.ToString(b.Name),
		})
	}
	table.Print(writer)

	return nil
}

func (p *BucketPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.buckets)
}

func (p *BucketPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.buckets)
}
