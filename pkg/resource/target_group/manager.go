package target_group

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun            bool
	elbClientv2       *aws.Elasticloadbalancingv2Client
	targetGroupGetter *Getter
}

func (m *Manager) Init() {
	if m.elbClientv2 == nil {
		m.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
	m.targetGroupGetter = NewGetter(m.elbClientv2)
}

func (m *Manager) Create(options resource.Options) error {
	tgOptions, ok := options.(*TargeGroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to TargeGroupOptions")
	}

	if tgOptions.Cluster != nil {
		tgOptions.VpcId = awssdk.ToString(tgOptions.Cluster.ResourcesVpcConfig.VpcId)
	}

	if m.DryRun {
		return m.dryRun(tgOptions, tgOptions.VpcId)
	}

	err := m.elbClientv2.CreateTargetGroup(tgOptions.Name, 1, tgOptions.Protocol, tgOptions.TargetType, tgOptions.VpcId)
	if err != nil {
		return err
	}
	fmt.Printf("Target Group %q created successfully\n", tgOptions.Name)

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	name := options.Common().Name

	tg, err := m.targetGroupGetter.GetTargetGroupByName(name)
	if err != nil {
		return err
	}

	if err := m.elbClientv2.DeleteTargetGroup(awssdk.ToString(tg.TargetGroupArn)); err != nil {
		return err
	}
	fmt.Printf("Target Group %q deleted\n", name)

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(options *TargeGroupOptions, vpcId string) error {
	fmt.Println("\nTarget Group Resource Manager Dry Run:")

	fmt.Printf("Elastic Load Balancing API Call %q with request parameters:\n", "CreateTargetGroup")
	fmt.Printf("Name: %q\n", options.Name)
	fmt.Printf("Port: %q\n", "1")
	fmt.Printf("Protocol: %q\n", options.Protocol)
	fmt.Printf("TargetType: %q\n", options.TargetType)
	fmt.Printf("VpcId: %q\n\n", vpcId)

	return nil
}
