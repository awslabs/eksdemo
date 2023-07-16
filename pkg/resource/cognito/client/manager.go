package client

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
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
		return fmt.Errorf("internal error, unable to cast options to client.Options")
	}

	if m.DryRun {
		return m.dryRun(options)
	}

	fmt.Printf("Creating App Client %q for User Pool Id %q...", options.ClientName, options.UserPoolID)
	appClient, err := m.cognitoClient.CreateUserPoolClient(
		options.OAuthScopes,
		options.CallbackUrls,
		options.ClientName,
		options.UserPoolID,
	)
	if err != nil {
		return aws.FormatError(err)
	}
	fmt.Printf("done\nCreated Cognito App Client Id: %s\n", awssdk.ToString(appClient.ClientId))

	return nil
}

func (m *Manager) Delete(_ resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(_ resource.Options, _ *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *Options) error {
	fmt.Printf("\nCognito App Client Resource Manager Dry Run:\n")
	fmt.Printf("Cognito User Pool API Call %q with request parameters:\n", "CreateUserPoolClient")
	fmt.Printf("AllowedOAuthFlows: %q\n", []types.OAuthFlowType{types.OAuthFlowTypeCode})
	fmt.Printf("AllowedOAuthScopes: %q\n", options.OAuthScopes)
	fmt.Printf("CallbackURLs: %q\n", options.CallbackUrls)
	fmt.Printf("ClientName: %q\n", options.ClientName)
	fmt.Printf("GenerateSecret: %s\n", "true")
	fmt.Printf("SupportedIdentityProviders: %q\n", []string{"COGNITO"})
	fmt.Printf("UserPoolId: %q\n", options.UserPoolID)
	return nil
}
