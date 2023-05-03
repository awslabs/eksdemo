package dns_record

import (
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

const MAX_COMBINED_NAME_AND_RECORD_LENGTH int = 90

type RecordSetPrinter struct {
	recordSets        []types.ResourceRecordSet
	longestNameLength int
}

func NewPrinter(recordSets []types.ResourceRecordSet) *RecordSetPrinter {
	return &RecordSetPrinter{recordSets: recordSets}
}

func (p *RecordSetPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Value"})

	for _, rs := range p.recordSets {
		if l := len(aws.ToString(rs.Name)); l > p.longestNameLength {
			p.longestNameLength = l
		}
	}

	for _, rs := range p.recordSets {
		records := ""
		// ALIAS records
		if rs.AliasTarget != nil {
			records = p.limitLength(aws.ToString(rs.AliasTarget.DNSName))
		} else if rs.Type == types.RRTypeSoa {
			// SOA records: split the MNAME, RNAME and rest into separate lines
			transform := strings.Replace(aws.ToString(rs.ResourceRecords[0].Value), " ", ",", 2)
			for i, rec := range strings.Split(transform, ",") {
				if i == 0 {
					records = p.limitLength(rec)
				} else {
					records += "\n" + p.limitLength(rec)
				}
			}
		} else {
			// All other records
			for i, rec := range rs.ResourceRecords {
				if i == 0 {
					records = p.limitLength(aws.ToString(rec.Value))
				} else {
					records += "\n" + p.limitLength(aws.ToString(rec.Value))
				}
			}
		}

		table.AppendRow([]string{
			strings.TrimSuffix(aws.ToString(rs.Name), "."),
			string(rs.Type),
			records,
		})
	}

	table.Print(writer)

	return nil
}

func (p *RecordSetPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.recordSets)
}

func (p *RecordSetPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.recordSets)
}

func (p *RecordSetPrinter) limitLength(record string) string {
	if len(record) > MAX_COMBINED_NAME_AND_RECORD_LENGTH-p.longestNameLength {
		record = record[:MAX_COMBINED_NAME_AND_RECORD_LENGTH-p.longestNameLength-3] + "..."
	}
	return record
}
