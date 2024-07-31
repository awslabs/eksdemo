package amp_workspace

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun           bool
	ampClient        *aws.AMPClient
	prometheusGetter *Getter
}

func (m *Manager) Init() {
	if m.ampClient == nil {
		m.ampClient = aws.NewAMPClient()
	}
	m.prometheusGetter = NewGetter(m.ampClient)
}

func (m *Manager) Create(options resource.Options) error {
	ampOptions, ok := options.(*AmpWorkspaceOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	_, err := m.prometheusGetter.GetAmpByAlias(ampOptions.Alias)
	// Return if the Workspace already exists
	if err == nil {
		fmt.Printf("AMP Workspace Alias %q already exists\n", ampOptions.Alias)
		return nil
	}

	// Return the error if it's anything other than resource not found
	if _, ok := err.(resource.NotFoundError); !ok {
		return err
	}

	if m.DryRun {
		return m.dryRun(ampOptions)
	}

	fmt.Printf("Creating AMP Workspace Alias: %s...", ampOptions.Alias)
	result, err := m.ampClient.CreateWorkspace(ampOptions.Alias)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated AMP Workspace Id: %s\n", awssdk.ToString(result.WorkspaceId))

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	ampOptions, ok := options.(*AmpWorkspaceOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpOptions")
	}

	id := options.Common().Id

	if id == "" {
		amp, err := m.prometheusGetter.GetAmpByAlias(ampOptions.Alias)
		if err != nil {
			if _, ok := err.(resource.NotFoundError); ok {
				fmt.Printf("AMP Workspace Alias %q does not exist\n", ampOptions.Alias)
				return nil
			}
			return err
		}
		id = awssdk.ToString(amp.Workspace.WorkspaceId)
	}

	err := m.ampClient.DeleteWorkspace(id)
	if err != nil {
		return aws.FormatError(err)
	}
	fmt.Printf("AMP Workspace Id %q deleting...\n", id)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *AmpWorkspaceOptions) error {
	fmt.Printf("\nAMP Resource Manager Dry Run:\n")
	fmt.Printf("Amazon Managed Service for Prometheus API Call %q with request parameters:\n", "CreateWorkspace")
	fmt.Printf("alias: %q\n", options.Alias)
	return nil
}
