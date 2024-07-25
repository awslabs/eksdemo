package core

import (
	"context"
	"fmt"

	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Check() *resource.Resource {
	return &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "check-for-crossplane-core",
		},

		Manager: &CheckResourceManager{
			Name:     "Crossplane Core",
			Message:  "please install \"crossplane-core\" first",
			Group:    "pkg.crossplane.io",
			Resource: "providers",
			Version:  "v1",
		},
	}
}

type CheckResourceManager struct {
	DryRun  bool
	Name    string
	Message string

	Group    string
	Resource string
	Version  string
	resource.EmptyInit
}

func (m *CheckResourceManager) Create(options resource.Options) error {
	if m.DryRun {
		fmt.Println("\nCheck Resource Manager Dry Run:")
		fmt.Println("Will check for: ", m.Name)
		fmt.Printf("Check performs a List operaton on Group: %q, Resource: %q, Version: %q\n", m.Group, m.Resource, m.Version)
		return nil
	}

	client, err := kubernetes.DynamicClient(options.Common().KubeContext)
	if err != nil {
		return err
	}

	providerConfig := schema.GroupVersionResource{
		Group:    m.Group,
		Resource: m.Resource,
		Version:  m.Version,
	}

	_, err = client.Resource(providerConfig).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf(m.Message)
	}

	fmt.Printf("%q found\n", m.Name)

	return nil
}

func (m *CheckResourceManager) Delete(_ resource.Options) error {
	return nil
}

func (m *CheckResourceManager) SetDryRun() {
	m.DryRun = true
}

func (m *CheckResourceManager) Update(_ resource.Options, _ *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
