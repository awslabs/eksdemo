package kyverno

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewRequireRequestsPolicy() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "require-requests",
			Description: "Require Resource Requests Policy",
		},

		Options: &resource.CommonOptions{
			Name: "kyverno-policy-require-requests",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: requireRequestsYamlTemplate,
			},
		},
	}
}

// https://github.com/dorukozturk/eks-best-practices-policies/blob/main/security/check-container-requests-limits.yaml
// https://aws.github.io/aws-eks-best-practices/security/docs/pods/#set-requests-and-limits-for-each-container-to-avoid-resource-contention-and-dos-attacks

const requireRequestsYamlTemplate = `---
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-requests-limits   
spec:
  validationFailureAction: audit
  background: true
  rules:
  - name: validate-resources
    match:
      resources:
        kinds:
        - Pod
    exclude:
      resources:
        kinds:
        - Pod
        namespaces:
        - kube-system
    validate:
      message: "CPU and memory resource requests and limits are required."
      pattern:
        spec:
          containers:
          - resources:
              requests:
                memory: "?*"
                cpu: "?*"
              limits:
                memory: "?*"
`
