package manifest

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

type ResourceManager struct {
	template.Template
	DryRun bool
	resource.EmptyInit
}

func (m *ResourceManager) Create(options resource.Options) error {
	manifest, err := m.Render(options)
	if err != nil {
		return err
	}

	if m.DryRun {
		fmt.Println("\nKubernetes Resource Manager Dry Run:")
		fmt.Println(manifest)
		return nil
	}

	return kubernetes.CreateResources(options.Common().KubeContext, manifest)
}

func (m *ResourceManager) Delete(options resource.Options) error {
	return fmt.Errorf("feature not yet implemented")
}

func (m *ResourceManager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not yet implemented")
}

func (m *ResourceManager) SetDryRun() {
	m.DryRun = true
}
