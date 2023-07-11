package userpool

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun         bool
	cognitoClient  *aws.CognitoUserPoolClient
	userPoolGetter *Getter
}

func (m *Manager) Init() {
	if m.cognitoClient == nil {
		m.cognitoClient = aws.NewCognitoUserPoolClient()
	}
	m.userPoolGetter = NewGetter(m.cognitoClient)
}

func (m *Manager) Create(options resource.Options) error {
	o, ok := options.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to userpool.Options")
	}

	_, err := m.userPoolGetter.GetUserPoolByName(o.UserPoolName)
	// Return if the User Pool already exists
	if err == nil {
		fmt.Printf("Cognito User Pool with name %q already exists\n", o.UserPoolName)
		return nil
	}

	// Return the error if it's anything other than resource not found
	var rnfe *resource.NotFoundByNameError
	if !errors.As(err, &rnfe) {
		return err
	}

	if m.DryRun {
		return m.dryRun(o)
	}

	fmt.Printf("Creating User Pool: %s...", o.UserPoolName)
	result, err := m.cognitoClient.CreateUserPool(o.UserPoolName)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated User Pool Id: %s\n", awssdk.ToString(result.Id))

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
	fmt.Printf("\nCognito User Pool Resource Manager Dry Run:\n")
	fmt.Printf("Cognito User Pool API Call %q with request parameters:\n", "CreateUserPool")
	fmt.Printf("PoolName: %q\n", options.UserPoolName)
	return nil
}
