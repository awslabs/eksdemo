package kms_key

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type KmsKeyPrinter struct {
	keys []*KMSKey
}

func NewPrinter(keys []*KMSKey) *KmsKeyPrinter {
	return &KmsKeyPrinter{keys}
}

func (p *KmsKeyPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Alias", "Id", "Status", "Key Spec"})

	for _, k := range p.keys {
		age := durafmt.ParseShort(time.Since(aws.ToTime(k.Key.CreationDate)))
		keyId := aws.ToString(k.Key.KeyId)
		alias := "-"

		if len(k.Aliases) > 0 {
			alias = strings.TrimPrefix(aws.ToString(k.Aliases[0].AliasName), "alias/")
		}
		if len(k.Aliases) > 1 {
			alias += fmt.Sprintf(" (+%d more)", len(k.Aliases)-1)
		}

		table.AppendRow([]string{
			age.String(),
			alias,
			keyId,
			string(k.Key.KeyState),
			string(k.Key.KeySpec),
		})
	}

	table.Print(writer)
	return nil
}

func (p *KmsKeyPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.keys)
}

func (p *KmsKeyPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.keys)
}
