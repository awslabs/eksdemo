package key

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun    bool
	kmsClient *aws.KMSClient
	kmsGetter *Getter
}

func (m *Manager) Init() {
	if m.kmsClient == nil {
		m.kmsClient = aws.NewKMSClient()
	}
	m.kmsGetter = NewGetter(m.kmsClient)
}

func (m *Manager) Create(options resource.Options) error {
	alias := options.Common().Name

	_, err := m.kmsGetter.GetByAlias(alias)

	// Return if the KMS alias already exists
	if err == nil {
		fmt.Printf("KMS Key with alias %q already exists\n", alias)
		return nil
	}

	// Return the error if it's anything other than resource not found
	var notFoundErr *resource.NotFoundByError
	if !errors.As(err, &notFoundErr) {
		return err
	}

	fullAliasName := fmt.Sprintf("alias/%s", alias)

	if m.DryRun {
		return m.dryRun(fullAliasName)
	}

	fmt.Printf("Creating KMS Key with Alias %q...", alias)

	keyMeta, err := m.kmsClient.CreateKey()
	if err != nil {
		return err
	}

	keyID := awssdk.ToString(keyMeta.KeyId)

	err = m.kmsClient.CreateAlias(fullAliasName, keyID)
	if err != nil {
		return fmt.Errorf("failed to create alias for key %q: %w", keyID, err)
	}
	fmt.Printf("done\nCreated KMS Key Id: %s\n", keyID)

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

func (m *Manager) dryRun(aliasName string) error {
	fmt.Printf("\nKMS Key Manager Dry Run:\n")
	fmt.Printf("KMS API Call %q with no request parameters\n", "CreateKey")
	fmt.Printf("KMS API Call %q with parameters:\n", "CreateAlias")
	fmt.Printf("\tAliasName: %q\n", aliasName)
	fmt.Printf("\tTargetKeyId: <Key Id returned from CreateKey call>\n")
	return nil
}
