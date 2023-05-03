package log_group

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun               bool
	cloudwatchlogsClient *aws.CloudwatchlogsClient
}

func (m *Manager) Init() {
	if m.cloudwatchlogsClient == nil {
		m.cloudwatchlogsClient = aws.NewCloudwatchlogsClient()
	}
}

func (m *Manager) Create(options resource.Options) error {
	_, err := m.cloudwatchlogsClient.CreateLogGroup(options.Common().Name)
	if err != nil {
		return aws.FormatError(err)
	}

	fmt.Printf("log-group %q created\n", options.Common().Name)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	logGroupName := options.Common().Name

	if err := m.cloudwatchlogsClient.DeleteLogGroup(logGroupName); err != nil {
		var rnfe *types.ResourceNotFoundException
		if errors.As(err, &rnfe) {
			return resource.NotFoundError(fmt.Sprintf("log-group %q does not exist", logGroupName))
		}
		return err
	}
	fmt.Printf("log-group %q deleted\n", logGroupName)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
