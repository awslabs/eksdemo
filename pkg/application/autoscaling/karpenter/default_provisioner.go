package karpenter

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type ProvisionerOptions struct {
	resource.CommonOptions
	*KarpenterOptions
}

func karpenterDefaultProvisioner(o *KarpenterOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &ProvisionerOptions{
			CommonOptions: resource.CommonOptions{
				Name: "karpenter-default-provisioner",
			},
			KarpenterOptions: o,
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

const yamlTemplate = `---
apiVersion: karpenter.sh/v1alpha5
kind: Provisioner
metadata:
  name: default
spec:
  requirements:
    - key: karpenter.sh/capacity-type
      operator: In
      values: ["on-demand", "spot"]
  limits:
    resources:
      cpu: 1000
  providerRef:
    name: default
{{- if .TTLSecondsAfterEmpty }}
  ttlSecondsAfterEmpty: {{ .TTLSecondsAfterEmpty }}
{{- else }}
  consolidation:
    enabled: true
{{- end }}
---
apiVersion: karpenter.k8s.aws/v1alpha1
kind: AWSNodeTemplate
metadata:
  name: default
spec:
  amiFamily: {{ .AMIFamily }}
  subnetSelector:
    Name: eksctl-{{ .ClusterName }}-cluster/SubnetPrivate*
  securityGroupSelector:
    aws:eks:cluster-name: {{ .ClusterName }}
`
