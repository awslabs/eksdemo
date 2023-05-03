package addon

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type AddonVersionPrinter struct {
	addonInfos []types.AddonInfo
}

func NewVersionPrinter(addonInfos []types.AddonInfo) *AddonVersionPrinter {
	return &AddonVersionPrinter{addonInfos}
}

func (p *AddonVersionPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Version", "Restrictions"})

	for _, addonInfo := range p.addonInfos {
		name := aws.ToString(addonInfo.AddonName)

		for _, av := range addonInfo.AddonVersions {
			isDefault := ""
			restrictions := "-"

			if len(av.Compatibilities) > 0 {
				if av.Compatibilities[0].DefaultVersion {
					isDefault = "*"
				}

				if len(av.Compatibilities[0].PlatformVersions) > 0 && av.Compatibilities[0].PlatformVersions[0] != "*" {
					restrictions = strings.Join((av.Compatibilities[0].PlatformVersions), ",")
				}
			}

			table.AppendRow([]string{
				name,
				aws.ToString(av.AddonVersion) + isDefault,
				restrictions,
			})
		}

	}

	table.Print(writer)
	if len(p.addonInfos) > 0 {
		fmt.Println("* Indicates default version")
	}

	return nil
}

func (p *AddonVersionPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.addonInfos)
}

func (p *AddonVersionPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.addonInfos)
}
