package amg_workspace

import (
	"errors"
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

type Manager struct {
	AssumeRolePolicyTemplate template.TextTemplate
	DryRun                   bool
	grafanaClient            *aws.GrafanaClient
	grafanaGetter            *Getter
	iamClient                *aws.IAMClient
}

func (m *Manager) Init() {
	if m.grafanaClient == nil {
		m.grafanaClient = aws.NewGrafanaClient()
	}
	m.grafanaGetter = NewGetter(m.grafanaClient)
	m.iamClient = aws.NewIAMClient()
}

func (m *Manager) Create(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	workspace, err := m.grafanaGetter.GetAmgByName(amgOptions.WorkspaceName)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if workspace != nil {
		fmt.Printf("AMG Workspace %q already exists\n", amgOptions.WorkspaceName)
		return nil
	}

	if m.DryRun {
		return m.dryRun(amgOptions)
	}

	role, err := m.createIamRole(amgOptions)
	if err != nil {
		return err
	}

	err = m.iamClient.PutRolePolicy(awssdk.ToString(role.RoleName), rolePolicName, rolePolicyDoc)
	if err != nil {
		return err
	}

	fmt.Printf("Creating AMG Workspace Name: %s...", amgOptions.WorkspaceName)
	result, err := m.grafanaClient.CreateWorkspace(amgOptions.WorkspaceName, amgOptions.Auth, awssdk.ToString(role.Arn))
	if err != nil {
		fmt.Println()
		return aws.FormatError(err)
	}

	fmt.Printf("done\nCreated AMG Workspace Id: %s\n", awssdk.ToString(result.Id))

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	amgOptions, ok := options.(*AmgOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmgOptions")
	}

	var amg *types.WorkspaceDescription
	var err error

	if options.Common().Id == "" {
		amg, err = m.grafanaGetter.GetAmgByName(amgOptions.WorkspaceName)
		if err != nil {
			if _, ok := err.(resource.NotFoundError); ok {
				fmt.Printf("AMG Workspace Name %q does not exist\n", amgOptions.WorkspaceName)
				return nil
			}
			return err
		}
	} else {
		amg, err = m.grafanaClient.DescribeWorkspace(options.Common().Id)
		if err != nil {
			return aws.FormatError(err)
		}
	}

	err = m.deleteIamRole(awssdk.ToString(amg.WorkspaceRoleArn))
	if err != nil {
		return err
	}

	id := awssdk.ToString(amg.Id)

	err = m.grafanaClient.DeleteWorkspace(id)
	if err != nil {
		return err
	}
	fmt.Printf("AMG Workspace Id %q deleting...\n", id)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) createIamRole(options *AmgOptions) (*iamtypes.Role, error) {
	assumeRolePolicy, err := m.AssumeRolePolicyTemplate.Render(options)
	if err != nil {
		return nil, err
	}

	roleName := options.iamRoleName()

	role, err := m.iamClient.CreateRole(assumeRolePolicy, roleName, "/service-role/")
	if err != nil {
		var eaee *iamtypes.EntityAlreadyExistsException
		if errors.As(err, &eaee) {
			fmt.Printf("IAM Role %q already exists\n", roleName)
			return m.iamClient.GetRole(roleName)
		}
		return nil, err
	}

	fmt.Printf("Created IAM Role: %s\n", awssdk.ToString(role.RoleName))

	return role, nil
}

func (m *Manager) deleteIamRole(roleArn string) error {
	roleName := roleArn[strings.LastIndex(roleArn, "/")+1:]

	// Delete inline policies before deleting role
	inlinePolicyNames, err := m.iamClient.ListRolePolicies(roleName)
	if err != nil {
		// Return an error if it's anything other than IAM role doesn't exist
		var nsee *iamtypes.NoSuchEntityException
		if errors.As(err, &nsee) {
			return nil
		}
		return err
	}

	for _, policyName := range inlinePolicyNames {
		err := m.iamClient.DeleteRolePolicy(roleName, policyName)
		if err != nil {
			return err
		}
	}

	// Remove managed policies before deleting role
	mgdPolicies, err := m.iamClient.ListAttachedRolePolicies(roleName)
	if err != nil {
		return err
	}

	for _, policy := range mgdPolicies {
		err := m.iamClient.DetachRolePolicy(roleName, awssdk.ToString(policy.PolicyArn))
		if err != nil {
			return err
		}
	}

	return m.iamClient.DeleteRole(roleName)
}

func (m *Manager) dryRun(options *AmgOptions) error {
	fmt.Println("\nAMG Resource Manager Dry Run:")

	fmt.Printf("Amazon Managed Grafana API Call %q with request parameters:\n", "CreateWorkspace")
	fmt.Printf("AccountAccessType: %q\n", types.AccountAccessTypeCurrentAccount)
	fmt.Printf("AuthenticationProviders: %q\n", options.Auth)

	fmt.Printf("PermissionType: %q\n", types.PermissionTypeServiceManaged)
	fmt.Printf("WorkspaceDataSources: %q\n", []types.DataSourceType{types.DataSourceTypePrometheus})
	fmt.Printf("WorkspaceName: %q\n", options.WorkspaceName)
	fmt.Printf("WorkspaceRoleArn: %q\n", "<role-to-be-created>")

	return nil
}
