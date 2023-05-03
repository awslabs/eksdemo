package fargate_profile

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type FargateProfilePrinter struct {
	profiles []*types.FargateProfile
}

func NewPrinter(profiles []*types.FargateProfile) *FargateProfilePrinter {
	return &FargateProfilePrinter{profiles}
}

func (p *FargateProfilePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Status", "Name", "Selectors"})

	for _, profile := range p.profiles {
		age := durafmt.ParseShort(time.Since(aws.ToTime(profile.CreatedAt)))
		name := aws.ToString(profile.FargateProfileName)

		selectors := make([]string, 0, len(profile.Selectors))
		for _, s := range profile.Selectors {
			selectorYaml, _ := json.MarshalIndent(s, "", "")
			selectors = append(selectors, string(selectorYaml))
		}

		table.AppendRow([]string{
			age.String(),
			string(profile.Status),
			name,
			strings.Join(selectors, ","),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *FargateProfilePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.profiles)
}

func (p *FargateProfilePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.profiles)
}
