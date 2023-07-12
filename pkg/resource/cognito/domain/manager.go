package domain

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun        bool
	cognitoClient *aws.CognitoUserPoolClient
}

func (m *Manager) Init() {
	if m.cognitoClient == nil {
		m.cognitoClient = aws.NewCognitoUserPoolClient()
	}
}

func (m *Manager) Create(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to domain.Options")
	}

	if m.DryRun {
		return m.dryRun(options)
	}

	fmt.Printf("Creating Domain %q for User Pool Id %q...", options.DomainName, options.UserPoolID)
	_, err := m.cognitoClient.CreateUserPoolDomain(options.DomainName, options.UserPoolID)
	if err != nil {
		return aws.FormatError(err)
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) Delete(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to domain.Options")
	}

	err := m.cognitoClient.DeleteUserPoolDomain(options.DomainName, options.UserPoolID)
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}
	fmt.Printf("Cognito Domain %q deleted\n", options.DomainName)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(_ resource.Options, _ *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *Options) error {
	fmt.Printf("\nCognito Domain Resource Manager Dry Run:\n")
	fmt.Printf("Cognito User Pool API Call %q with request parameters:\n", "CreateUserPooDomain")
	fmt.Printf("Domain: %q\n", options.DomainName)
	fmt.Printf("UserPoolId: %q\n", options.UserPoolID)
	return nil
}
