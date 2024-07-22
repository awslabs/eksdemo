package core

import (
	"context"
	"fmt"
	"time"

	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

func waitForProviderCRD() *resource.Resource {
	return &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "wait-for-provider-crd",
		},

		Manager: &WaitManager{
			CRD:         "Provider",
			Group:       "pkg.crossplane.io",
			Resource:    "providers",
			Version:     "v1",
			WaitSeconds: 200,
		},
	}
}

func waitForProviderConfigCRD() *resource.Resource {
	return &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "wait-for-provider-config-crd",
		},

		Manager: &WaitManager{
			CRD:         "ProviderConfig",
			Group:       "aws.upbound.io",
			Resource:    "providerconfigs",
			Version:     "v1beta1",
			WaitSeconds: 200,
		},
	}
}

type WaitManager struct {
	CRD         string
	DryRun      bool
	Group       string
	Resource    string
	Version     string
	WaitSeconds int
	resource.EmptyInit
}

func (m *WaitManager) Create(options resource.Options) error {
	if m.DryRun {
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

	fiveMinRetry := wait.Backoff{
		Steps:    m.WaitSeconds / 2,
		Duration: 2000 * time.Millisecond,
	}

	fmt.Printf("Waiting for %q CRD to be created", m.CRD)

	err = retry.OnError(fiveMinRetry, errors.IsNotFound, func() (err error) {
		_, err = client.Resource(providerConfig).List(context.Background(), metav1.ListOptions{})
		fmt.Printf(".")
		return
	})
	fmt.Println()

	if err != nil {
		return fmt.Errorf("timed out waiting for CRD to be created")
	}

	return nil
}

func (m *WaitManager) Delete(_ resource.Options) error {
	return fmt.Errorf("feature not supported")

}

func (m *WaitManager) SetDryRun() {
	m.DryRun = true
}

func (m *WaitManager) Update(_ resource.Options, _ *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
