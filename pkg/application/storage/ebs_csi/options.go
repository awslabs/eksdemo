package ebs_csi

import (
	"context"
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EbsCsiOptions struct {
	application.ApplicationOptions

	NoStorageClasses bool
	Replicas         int
}

const IsDefaultStorageClassAnnotation = "storageclass.kubernetes.io/is-default-class"

func newOptions() (options *EbsCsiOptions, flags cmd.Flags) {
	options = &EbsCsiOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "kube-system",
			ServiceAccount: "ebs-csi-controller-sa",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2.19.0",
				Latest:        "v1.19.0",
				PreviousChart: "2.17.1",
				Previous:      "v1.16.1",
			},
		},
		NoStorageClasses: false,
		Replicas:         1,
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-storageclasses",
				Description: "don't create the gp3 Storage Classes",
			},
			Option: &options.NoStorageClasses,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}

func (o *EbsCsiOptions) PreInstall() error {
	if o.NoStorageClasses {
		return nil
	}

	if o.DryRun {
		fmt.Println("\nPreInstall Dry Run:")
		fmt.Println("Mark the current default StorageClass as non-default")
		return nil
	}

	fmt.Println("Checking for default StorageClass")
	k8sclient, err := kubernetes.Client(o.KubeContext())
	if err != nil {
		return err
	}

	scList, err := k8sclient.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, sc := range scList.Items {
		if sc.Annotations[IsDefaultStorageClassAnnotation] != "true" {
			continue
		}
		fmt.Printf("Marking StorageClass %q as non-default...", sc.Name)

		sc.Annotations[IsDefaultStorageClassAnnotation] = "false"
		k8sclient.StorageV1().StorageClasses().Update(context.Background(), &sc, metav1.UpdateOptions{})

		fmt.Println("done")
	}

	return nil
}
