package organization

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	organizationsClient *aws.OrganizationsClient
}

func (m *Manager) Init() {
	if m.organizationsClient == nil {
		m.organizationsClient = aws.NewOrganizationsClient()
	}
}

func (m *Manager) Create(options resource.Options) error {
	result, err := m.organizationsClient.CreateOrganization()
	if err != nil {
		return err
	}
	fmt.Printf("Created AWS Organization: %s\n", awssdk.ToString(result.Id))

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	err := m.organizationsClient.DeleteOrganization()
	if err != nil {
		return err
	}
	fmt.Println("AWS Organization deleted")

	return nil
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("update not supported")
}

func (m *Manager) SetDryRun() {}
