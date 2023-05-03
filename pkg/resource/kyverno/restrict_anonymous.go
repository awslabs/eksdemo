package kyverno

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewRestrictAnonymousPolicy() *resource.Resource {
	return &resource.Resource{
		Command: cmd.Command{
			Name:        "restrict-anonymous",
			Description: "Restrict Anonymous Access Policy",
		},

		Options: &resource.CommonOptions{
			Name: "kyverno-policy-restrict-anonymous",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: restrictAnonymousYamlTemplate,
			},
		},
	}
}

// https://github.com/dorukozturk/eks-best-practices-policies/blob/main/security/restrict-anonymous-access.yaml
// https://aws.github.io/aws-eks-best-practices/security/docs/iam/#review-and-revoke-unnecessary-anonymous-access

const restrictAnonymousYamlTemplate = `---
apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: restrict-anonymous-access
spec:
  validationFailureAction: audit
  background: true
  rules:
    - name: unauthenticated-access
      match:
        any:
        - resources:
            kinds:
              - RoleBinding
              - ClusterRoleBinding
      exclude:
        resources:
          kinds:
          - RoleBinding
          - ClusterRoleBinding
          namespaces:
          - kube-system
      validate:
        message: "Don't bind roles or clusterroles to system:unauthenticated group."
        pattern:
          subjects:
            name: "!system:unauthenticated"
    - name: anonymous-access
      match:
        any:
        - resources:
            kinds:
              - RoleBinding
              - ClusterRoleBinding
      exclude:
        resources:
          kinds:
          - RoleBinding
          - ClusterRoleBinding
          namespaces:
          - kube-system
      validate:
        message: "Don't bind roles or clusterroles to system:anonymous group."
        pattern:
          subjects:
            name: "!system:anonymous"
`
