package iam_role

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type IamRolePrinter struct {
	roles    []types.Role
	lastUsed bool
}

func NewPrinter(roles []types.Role, lastUsed bool) *IamRolePrinter {
	return &IamRolePrinter{roles, lastUsed}
}

func (p *IamRolePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()

	header := []string{"Age", "Role"}
	if p.lastUsed {
		header = append(header, "Last Used")
	}

	table.SetHeader(header)

	for _, r := range p.roles {
		age := durafmt.ParseShort(time.Since(aws.ToTime(r.CreateDate)))
		rlu := r.RoleLastUsed

		row := []string{
			age.String(),
			aws.ToString(r.RoleName),
		}

		if p.lastUsed {
			var lastUsed string

			if rlu != nil && rlu.LastUsedDate != nil {
				lastUsed = durafmt.ParseShort(time.Since(aws.ToTime(rlu.LastUsedDate))).String()
			} else {
				lastUsed = "-"
			}
			row = append(row, lastUsed)
		}

		table.AppendRow(row)
	}

	table.Print(writer)

	return nil
}

func (p *IamRolePrinter) PrintJSON(writer io.Writer) error {
	if err := p.decodeAssumeRolePolicyDocuments(); err != nil {
		return err
	}
	return printer.EncodeJSON(writer, p.roles)
}

func (p *IamRolePrinter) PrintYAML(writer io.Writer) error {
	if err := p.decodeAssumeRolePolicyDocuments(); err != nil {
		return err
	}
	return printer.EncodeYAML(writer, p.roles)
}

func (p *IamRolePrinter) decodeAssumeRolePolicyDocuments() error {
	for i, r := range p.roles {
		decodedValue, err := url.QueryUnescape(aws.ToString(r.AssumeRolePolicyDocument))
		if err != nil {
			return fmt.Errorf("unable to decode AssumeRolePolicyDocument for role %q: %w", aws.ToString(r.RoleName), err)

		}
		p.roles[i].AssumeRolePolicyDocument = aws.String(decodedValue)
	}
	return nil
}
