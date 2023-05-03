package iam_policy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type RolePolicies struct {
	InlinePolicies  []InlinePolicy
	ManagedPolicies []ManagedPolicy
}

type InlinePolicy struct {
	Name           string
	PolicyDocument string
}

type ManagedPolicy struct {
	Policy        *types.Policy
	PolicyVersion *types.PolicyVersion
}

type RolePolicyPrinter struct {
	*RolePolicies
}

func NewRolePolicyPrinter(policies *RolePolicies) *RolePolicyPrinter {
	return &RolePolicyPrinter{policies}
}

func (p *RolePolicyPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Name", "Type", "Description"})

	for _, p := range p.RolePolicies.InlinePolicies {
		table.AppendRow([]string{
			p.Name,
			"Inline",
			"",
		})
	}

	for _, p := range p.RolePolicies.ManagedPolicies {
		policyType := "failed to parse arn"

		arn, err := arn.Parse(aws.ToString(p.Policy.Arn))
		if err == nil {
			if arn.AccountID == "aws" {
				policyType = "AWS Mgd"
			} else {
				policyType = "Cust Mgd"
			}
		}

		table.AppendRow([]string{
			aws.ToString(p.Policy.PolicyName),
			policyType,
			aws.ToString(p.Policy.Description),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *RolePolicyPrinter) PrintJSON(writer io.Writer) error {
	if err := p.decodePolicyDocuments(); err != nil {
		return err
	}
	p.cleanPolicyDocuments()

	return printer.EncodeJSON(writer, p.RolePolicies)
}

func (p *RolePolicyPrinter) PrintYAML(writer io.Writer) error {
	if err := p.decodePolicyDocuments(); err != nil {
		return err
	}

	if err := p.prettyPrintPolicyDocuments(); err != nil {
		return err
	}

	return printer.EncodeYAML(writer, p.RolePolicies)
}

func (p *RolePolicyPrinter) cleanPolicyDocuments() {
	for i := range p.InlinePolicies {
		p.InlinePolicies[i].PolicyDocument = strings.Join(strings.Fields(p.InlinePolicies[i].PolicyDocument), " ")
	}

	for i := range p.ManagedPolicies {
		p.ManagedPolicies[i].PolicyVersion.Document =
			aws.String(strings.Join(strings.Fields(aws.ToString(p.ManagedPolicies[i].PolicyVersion.Document)), " "))
	}
}

func (p *RolePolicyPrinter) decodePolicyDocuments() error {
	for i, pol := range p.InlinePolicies {
		decodedValue, err := url.QueryUnescape(pol.PolicyDocument)
		if err != nil {
			return fmt.Errorf("unable to decode PolicyDocument for inline policy %q: %w", pol.Name, err)

		}
		p.InlinePolicies[i].PolicyDocument = decodedValue
	}

	for i, pol := range p.ManagedPolicies {
		decodedValue, err := url.QueryUnescape(aws.ToString(pol.PolicyVersion.Document))
		if err != nil {
			return fmt.Errorf("unable to decode PolicyDocument for managed policy %q: %w",
				aws.ToString(pol.Policy.Arn), err)

		}
		p.ManagedPolicies[i].PolicyVersion.Document = aws.String(decodedValue)
	}
	return nil
}

func (p *RolePolicyPrinter) prettyPrintPolicyDocuments() error {
	var prettyJSON bytes.Buffer

	for i, pol := range p.InlinePolicies {
		err := json.Indent(&prettyJSON, []byte(p.InlinePolicies[i].PolicyDocument), "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON for inline policy %q: %w", pol.Name, err)
		}

		p.InlinePolicies[i].PolicyDocument = prettyJSON.String()
	}

	for i, pol := range p.ManagedPolicies {
		err := json.Indent(&prettyJSON, []byte(aws.ToString(p.ManagedPolicies[i].PolicyVersion.Document)), "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON for managed policy %q: %w",
				aws.ToString(pol.Policy.Arn), err)
		}

		p.ManagedPolicies[i].PolicyVersion.Document = aws.String(prettyJSON.String())
	}
	return nil
}
