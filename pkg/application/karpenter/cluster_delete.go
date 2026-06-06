package karpenter

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/awslabs/eksdemo/pkg/kubernetes"
	api_errors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/dynamic"
	cmdwait "k8s.io/kubectl/pkg/cmd/wait"
)

func DeleteCustomResources(kubeContext string) error {
	client, err := kubernetes.DynamicClient(kubeContext)
	if err != nil {
		fmt.Printf("Warning: failed creating kubernetes dynamic client: %s\n", err)
		fmt.Println("Skipping cleaning up Karpenter resource.")
		return nil
	}

	if err := DeleteNodePool(client); err != nil {
		return err
	}

	return DeleteEC2NodeClass(client, kubeContext)
}

// This uses cli-runtime's Resource Builder and kubectl's waiter to simplify the delete and wait
// Ideally this should use the Dynamic Client and a custom waiter
func DeleteEC2NodeClass(client dynamic.Interface, kubeContext string) error {
	getter := genericclioptions.NewConfigFlags(true)
	getter.Context = &kubeContext

	restMapper, err := getter.ToRESTMapper()
	if err != nil {
		return fmt.Errorf("failed creating restMapper: %w", err)
	}

	gvk := "EC2NodeClass.v1.karpenter.k8s.aws"
	fullySpecifiedGVK, _ := schema.ParseKindArg(gvk)

	_, err = restMapper.RESTMapping(fullySpecifiedGVK.GroupKind(), fullySpecifiedGVK.Version)
	if err != nil {
		if meta.IsNoMatchError(err) {
			// EC2NodeClass kind doesn't exist, skip deletion
			return nil
		}
		return err
	}

	infos, err := resource.NewBuilder(getter).Unstructured().ResourceTypes(gvk).SelectAllParam(true).Flatten().Do().Infos()
	if err != nil {
		return err
	}

	for _, info := range infos {
		fmt.Printf("Deleting Karpenter EC2NodeClass %q\n", info.Name)

		_, err := resource.NewHelper(info.Client, info.Mapping).Delete(info.Namespace, info.Name)
		if err != nil {
			return fmt.Errorf("failed to delete EC2NodeClass %q: %w", info.Name, err)
		}
	}

	fmt.Println("Waiting up to 30 seconds for EC2NodeClasses to be deleted...")
	waitOptions := cmdwait.WaitOptions{
		ResourceFinder: genericclioptions.ResourceFinderForResult(resource.InfoListVisitor(infos)),
		DynamicClient:  client,
		Timeout:        30 * time.Second,

		Printer:     &printers.NamePrinter{ShortOutput: false, Operation: "deleted"},
		ConditionFn: cmdwait.IsDeleted,
		IOStreams:   genericiooptions.IOStreams{Out: os.Stdout, ErrOut: os.Stderr},
	}

	return waitOptions.RunWait()
}

func DeleteNodePool(client dynamic.Interface) error {
	nodePools := schema.GroupVersionResource{
		Group:    "karpenter.sh",
		Version:  "v1",
		Resource: "nodepools",
	}

	list, err := client.Resource(nodePools).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		if api_errors.IsNotFound(err) {
			// nodepools resource doesn't exist, skip deletion
			return nil
		}
		return err
	}

	for _, item := range list.Items {
		nodePoolName := item.GetName()

		fmt.Printf("Deleting Karpenter NodePool %q\n", nodePoolName)
		if err := client.Resource(nodePools).Delete(context.Background(), nodePoolName, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("failed to delete NodePool %q: %s\n", nodePoolName, err)
		}
	}
	return nil
}
