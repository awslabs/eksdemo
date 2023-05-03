package logs_insights

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/query"
)

const apiServerFilter = `filter @logStream like /^kube-apiserver-[0-9a-f]*$/`
const authenticatorFilter = `filter @logStream like /^authenticator-[0-9a-f]*$/`
const cloudControllerManagerFilter = `filter @logStream like /^cloud-controller-manager-[0-9a-f]*$/`
const controllerManagerFilter = `filter @logStream like /^kube-controller-manager-[0-9a-f]*$/`
const controlPlaneQuery = `filter @logStream not like /kube-apiserver-audit-[0-9a-f]*$/`
const schedulerFilter = `filter @logStream like /^kube-scheduler-[0-9a-f]*$/`

func NewApiServerQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = apiServerFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "apiserver",
			Description: "EKS API Server Logs Query",
			Aliases:     []string{"api"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewAuthenticatorQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = authenticatorFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "authenticator",
			Description: "EKS Authenticator Logs Query",
			Aliases:     []string{"auth"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewCloudControllerManagerQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = cloudControllerManagerFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "cloud-controller-manager",
			Description: "EKS Cloud Controller Manager Logs Query",
			Aliases:     []string{"cloud-controller", "ccm"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewControllerManagerQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = controllerManagerFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "controller-manager",
			Description: "EKS Controller Manager Logs Query",
			Aliases:     []string{"cm"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewControlPlaneQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = controlPlaneQuery

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "control-plane",
			Description: "EKS Control Plane Logs (excl. Audit) Query",
			Aliases:     []string{"cp"},
		},

		CreateFlags: createFlags,
		Manager:     &query.Manager{},
		Options:     options,
	}
}

func NewSchedulerQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = schedulerFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "scheduler",
			Description: "EKS Scheduler Logs Query",
			Aliases:     []string{"schedule", "sched"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}
