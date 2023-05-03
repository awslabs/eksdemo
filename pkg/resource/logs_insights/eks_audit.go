package logs_insights

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/query"
)

const auditFilter = `filter @logStream like /^kube-apiserver-audit-[0-9a-f]*$/`

const audit401Query = `
| filter responseStatus.code="401"`

const audit403Query = `
| filter responseStatus.code="403"`

func NewAuditQuery() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = auditFilter

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "audit",
			Description: "EKS Audit Logs Query",
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewAudit401Query() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = auditFilter + audit401Query

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "audit-401",
			Description: "EKS Audit Logs 401 Unauthorized Query",
			Aliases:     []string{"401"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}

func NewAudit403Query() *resource.Resource {
	options, createFlags := query.NewOptions()
	options.Query = auditFilter + audit403Query

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "audit-403",
			Description: "EKS Audit Logs 403 Forbidden Query",
			Aliases:     []string{"403"},
		},

		CreateFlags: createFlags,
		Options:     options,
		Manager:     &query.Manager{},
	}
}
