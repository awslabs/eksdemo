package client

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun          bool
	cognitoClient   *aws.CognitoUserPoolClient
	appClientGetter *Getter
}

func (m *Manager) Init() {
	if m.cognitoClient == nil {
		m.cognitoClient = aws.NewCognitoUserPoolClient()
	}
	m.appClientGetter = NewGetter(m.cognitoClient)
}

func (m *Manager) Create(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to client.Options")
	}

	if m.DryRun {
		return m.dryRun(options)
	}

	fmt.Printf("Creating App Client %q for User Pool Id %q...", options.AppClientName, options.UserPoolID)
	appClient, err := m.cognitoClient.CreateUserPoolClient(
		options.OAuthScopes,
		options.CallbackUrls,
		options.AppClientName,
		options.UserPoolID,
	)
	if err != nil {
		return aws.FormatError(err)
	}
	fmt.Printf("done\nCreated Cognito App Client Id: %s\n", awssdk.ToString(appClient.ClientId))

	return nil
}

func (m *Manager) Delete(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to client.Options")
	}

	id := options.AppClientID

	if id == "" {
		ac, err := m.appClientGetter.GetAppClientByName(options.AppClientName, options.UserPoolID)

		if err != nil {
			var rnfe *resource.NotFoundByNameError
			if errors.As(err, &rnfe) {
				fmt.Printf("Cognito App Client with name %q does not exist\n", options.AppClientName)
				return nil
			}
			return err
		}
		id = awssdk.ToString(ac.ClientId)
	}

	err := m.cognitoClient.DeleteUserPoolClient(id, options.UserPoolID)
	if err != nil {
		return aws.FormatError(err)
	}
	fmt.Printf("Cognito App Client Id %q deleted\n", id)

	return nil
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
	fmt.Printf("ClientName: %q\n", options.AppClientName)
	fmt.Printf("GenerateSecret: %s\n", "true")
	fmt.Printf("SupportedIdentityProviders: %q\n", []string{"COGNITO"})
	fmt.Printf("UserPoolId: %q\n", options.UserPoolID)
	return nil
}
