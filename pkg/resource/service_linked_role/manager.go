package service_linked_role

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/iam_role"

	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun     bool
	iamClient  *aws.IAMClient
	roleGetter *iam_role.Getter
}

func (m *Manager) Init() {
	if m.iamClient == nil {
		m.iamClient = aws.NewIAMClient()
	}
	m.roleGetter = iam_role.NewGetter(m.iamClient)
}

func (m *Manager) Create(options resource.Options) error {
	slrOptions, ok := options.(*ServiceLinkedRoleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ServiceLinkedRoleOptions")
	}

	_, err := m.roleGetter.GetRoleByName(slrOptions.RoleName)

	// Return if the SLR already exists
	if err == nil {
		fmt.Printf("Service Linked Role %q already exists\n", slrOptions.RoleName)
		return nil
	}

	// Return the error if it's anything other than resource not found
	if _, ok := err.(resource.NotFoundError); !ok {
		return err
	}

	if m.DryRun {
		return m.dryRun(slrOptions)
	}

	fmt.Printf("Creating Service Linked Role: %s...", slrOptions.RoleName)

	err = m.iamClient.CreateServiceLinkedRole(slrOptions.ServiceName)
	if err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	slrOptions, ok := options.(*ServiceLinkedRoleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ServiceLinkedRoleOptions")
	}

	fmt.Printf("Service Linked Role %q will NOT be deleted\n", slrOptions.RoleName)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *ServiceLinkedRoleOptions) error {
	fmt.Printf("\nService Linked Role Manager Dry Run:\n")
	fmt.Printf("IAM API Call %q with request parameters:\n", "CreateServiceLinkedRole")
	fmt.Printf("AWSServiceName: %q\n", options.ServiceName)
	return nil
}
