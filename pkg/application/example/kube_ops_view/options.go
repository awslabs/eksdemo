package kube_ops_view

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type KubeOpsViewOptions struct {
	application.ApplicationOptions

	Replicas int
}

func newOptions() (options *KubeOpsViewOptions, flags cmd.Flags) {
	options = &KubeOpsViewOptions{
		ApplicationOptions: application.ApplicationOptions{
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "kube-ops-view",
			DefaultVersion: &application.LatestPrevious{
				Latest:   "latest",
				Previous: "20.4.0",
			},
			DisableServiceAccountFlag: true,
		},
		Replicas: 1,
	}
	return
}
