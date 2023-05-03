package ecr_repository

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type RepositoryPrinter struct {
	repos []types.Repository
}

func NewPrinter(repos []types.Repository) *RepositoryPrinter {
	return &RepositoryPrinter{repos}
}

func (p *RepositoryPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Tags", "Scan", "Encryption"})

	caser := cases.Title(language.English)

	for _, repo := range p.repos {
		age := durafmt.ParseShort(time.Since(aws.ToTime(repo.CreatedAt)))
		scan := "Manual"
		if repo.ImageScanningConfiguration.ScanOnPush {
			scan = "Scan on push"
		}

		table.AppendRow([]string{
			age.String(),
			aws.ToString(repo.RepositoryName),
			caser.String(string(repo.ImageTagMutability)),
			scan,
			string(repo.EncryptionConfiguration.EncryptionType),
		})
	}

	table.Print(writer)

	return nil
}

func (p *RepositoryPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.repos)
}

func (p *RepositoryPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.repos)
}
