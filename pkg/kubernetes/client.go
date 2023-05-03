package kubernetes

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/util/retry"
)

func CreateResources(kubeContext, manifest string) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.Context = &kubeContext
	builder := resource.NewBuilder(getter).Unstructured().ContinueOnError()

	switch {
	case strings.Index(manifest, "http://") == 0 || strings.Index(manifest, "https://") == 0:
		url, err := url.Parse(manifest)
		if err != nil {
			return fmt.Errorf("the URL %q is not valid: %v", manifest, err)
		}
		builder.URL(3, url)
	default:
		manifest := bytes.NewBufferString(manifest)
		builder.Stream(manifest, "manifest")
	}

	infos, err := builder.Do().Infos()

	if err != nil {
		return err
	}

	// Initial use case for Retry is Karpenter. Helm wait doesn't wait for Webhooks
	// https://github.com/cert-manager/cert-manager/issues/2908
	tenSecondRetry := wait.Backoff{
		Steps:    5,
		Duration: 2 * time.Second,
	}

	for _, info := range infos {
		fmt.Printf("Creating %s %q", info.Object.GetObjectKind().GroupVersionKind().Kind, info.Name)
		if info.Namespace != "" {
			fmt.Printf(" in namespace %q", info.Namespace)
		}
		fmt.Println()

		err = retry.OnError(tenSecondRetry, errors.IsInternalError, func() (err error) {
			_, err = resource.NewHelper(info.Client, info.Mapping).Create(info.Namespace, true, info.Object)
			return
		})

		if err != nil {
			fmt.Printf("Warning: failed to create resource: %s\n", err)
		}
	}

	return nil
}

func DeleteResources(kubeContext, manifestYaml string) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.Context = &kubeContext

	manifest := bytes.NewBufferString(manifestYaml)
	infos, err := resource.NewBuilder(getter).Unstructured().Stream(manifest, "test").Do().Infos()

	if err != nil {
		return err
	}

	for _, info := range infos {
		fmt.Printf("Deleting %s: %s", info.Object.GetObjectKind().GroupVersionKind().Kind, info.Name)
		if info.Namespace != "" {
			fmt.Printf(" in namespace: %s", info.Namespace)
		}
		fmt.Println()

		obj, err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name)
		if err != nil {
			fmt.Printf("Warning: failed to delete resource: %s\n", err)
		}
		_ = obj
	}

	return nil
}
