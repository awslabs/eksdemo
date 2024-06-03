package instance

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
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

func (m *Manager) Create(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to instance.Options")
	}

	if m.DryRun {
		return m.dryRun(options)
	}

	_, err := m.ec2Client.RunInstances(
		options.AMI,
		options.InstanceType,
		options.KeyName,
		options.Common().Name,
		options.SubnetID,
		options.volumeDeviceName,
		options.Count,
		options.VolumeSize,
	)

	return aws.FormatError(err)
}

func (m *Manager) Delete(options resource.Options) (err error) {
	instanceId := options.Common().Name

	ec2, err := m.ec2Getter.GetInstanceById(instanceId)
	if err != nil {
		return err
	}

	if ec2.State.Name == types.InstanceStateNameTerminated {
		return fmt.Errorf("ec2-instance %q already terminated", instanceId)
	}

	if err := m.ec2Client.TerminateInstances(instanceId); err != nil {
		return err
	}
	fmt.Printf("ec2-instance %q terminating...\n", instanceId)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *Options) error {
	fmt.Println("EC2 Instance Resource Manager Dry Run:")

	fmt.Printf("EC2 API Call %q with request parameters:\n", "RunInstances")
	fmt.Printf("InstanceType: %q\n", options.InstanceType)
	fmt.Printf("ImageId: %q\n", options.AMI)
	if options.KeyName != "" {
		fmt.Printf("KeyName: %q\n", options.KeyName)
	}
	fmt.Printf("MaxCount: %d\n", options.Count)
	fmt.Printf("MinCount: %d\n", options.Count)
	if options.SubnetID != "" {
		fmt.Printf("SubnetId: %q\n", options.SubnetID)
	}
	fmt.Printf("TagSpecifications: %q\n", fmt.Sprintf("ResourceType=instance,Tags=[{Key=Name,Value=%s}]", options.Name))
	if options.VolumeSize > 0 {
		fmt.Printf("BlockDeviceMappings: %s\n", fmt.Sprintf("[ { \"DeviceName\": %q, \"Ebs\": { \"VolumeSize\": %d } } ]",
			options.volumeDeviceName, options.VolumeSize),
		)
	}

	return nil
}
