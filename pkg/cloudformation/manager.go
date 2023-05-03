package cloudformation

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

// eksdemo-<clusterName>-<resourceName>
const stackNameTemplate = "eksdemo-%s-%s"

type ResourceManager struct {
	Resource string

	Capabilities []types.Capability
	DryRun       bool
	Parameters   map[string]string
	Template     template.Template

	cloudformationClient *aws.CloudformationClient
	cloudformationGetter *cloudformation_stack.Getter
}

func (m *ResourceManager) Init() {
	if m.cloudformationClient == nil {
		m.cloudformationClient = aws.NewCloudformationClient()
	}
	m.cloudformationGetter = cloudformation_stack.NewGetter(m.cloudformationClient)
}

func (m *ResourceManager) Create(options resource.Options) error {
	cfTemplate, err := m.Template.Render(options)
	if err != nil {
		return err
	}

	stackName := fmt.Sprintf(stackNameTemplate, options.Common().ClusterName, options.Common().Name)

	stack, err := m.cloudformationGetter.GetStacks(stackName)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if len(stack) > 0 {
		fmt.Printf("CloudFormation stack %q already exists\n", stackName)
		return nil
	}

	if m.DryRun {
		fmt.Println("\nCloudFormation Resource Manager Dry Run:")
		fmt.Printf("Stack name %q template:\n", stackName)
		fmt.Println(cfTemplate)
		return nil
	}

	fmt.Printf("Creating CloudFormation stack %q (can take 1+ min)...", stackName)
	err = m.cloudformationClient.CreateStack(stackName, cfTemplate, m.Parameters, m.Capabilities)

	if err != nil {
		fmt.Println()
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *ResourceManager) Delete(options resource.Options) error {
	stackName := fmt.Sprintf(stackNameTemplate, options.Common().ClusterName, options.Common().Name)

	_, err := m.cloudformationGetter.GetStacks(stackName)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); ok {
			fmt.Printf("CloudFormation Stack %q does not exist\n", stackName)
			return nil
		}
		return err
	}

	fmt.Printf("Deleting CloudFormation Stack %q\n", stackName)

	return m.cloudformationClient.DeleteStack(stackName)
}

func (m *ResourceManager) SetDryRun() {
	m.DryRun = true
}

func (m *ResourceManager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
