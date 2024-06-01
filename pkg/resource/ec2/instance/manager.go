package instance

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun    bool
	ec2Client *aws.EC2Client
	ec2Getter *Getter
}

func (m *Manager) Init() {
	if m.ec2Client == nil {
		m.ec2Client = aws.NewEC2Client()
	}
	m.ec2Getter = NewGetter(m.ec2Client)
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) (err error) {
	instanceId := options.Common().Name

	ec2, err := m.ec2Getter.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if string(ec2.State.Name) == "terminated" {
		return fmt.Errorf("ec2-instance %q already terminated", instanceId)
	}

	if err := m.ec2Client.TerminateInstances(instanceId); err != nil {
		return err
	}
	fmt.Println("EC2 Instance terminating...")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
