package eksctl

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

const EksctlHeader = `---
apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: {{ .ClusterName }}
  region: {{ .Region }}
{{- if .KubernetesVersion }}
  version: {{ .KubernetesVersion | printf "%q" }}
{{- end }}
`

type ResourceManager struct {
	Resource       string
	ConfigTemplate template.Template
	DeleteFlags    template.Template
	ApproveCreate  bool
	ApproveDelete  bool
	DryRun         bool
	resource.EmptyInit
}

func (e *ResourceManager) Create(options resource.Options) error {
	eksctlConfig, err := e.ConfigTemplate.Render(options)
	if err != nil {
		return err
	}

	args := []string{
		"create",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveCreate {
		args = append(args, "--approve")
	}

	if e.DryRun {
		fmt.Println("\nEksctl Resource Manager Dry Run:")
		fmt.Println("eksctl " + strings.Join(args, " "))
		fmt.Println(eksctlConfig)
		return nil
	}

	return Command(args, eksctlConfig)
}

func (e *ResourceManager) Delete(options resource.Options) error {
	if e.DeleteFlags != nil {
		return e.DeleteWithFlags(options)
	}
	return e.DeleteWithConfigFile(options)
}

func (e *ResourceManager) DeleteWithFlags(options resource.Options) error {
	deleteCommand, err := e.DeleteFlags.Render(options)
	if err != nil {
		return err
	}

	args := append([]string{"delete", e.Resource}, strings.Split(deleteCommand, " ")...)

	return Command(args, "")
}

func (e *ResourceManager) DeleteWithConfigFile(options resource.Options) error {
	eksctlConfig, err := e.ConfigTemplate.Render(options)

	if err != nil {
		return err
	}

	args := []string{
		"delete",
		e.Resource,
		"-f",
		"-",
	}

	if e.ApproveDelete {
		args = append(args, "--approve")
	}

	return Command(args, eksctlConfig)
}

func (e *ResourceManager) SetDryRun() {
	e.DryRun = true
}

func (e *ResourceManager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
