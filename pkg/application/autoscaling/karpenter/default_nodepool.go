package karpenter

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type NodePoolOptions struct {
	resource.CommonOptions
	*KarpenterOptions
}

func karpenterDefaultNodePool(o *KarpenterOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &NodePoolOptions{
			CommonOptions: resource.CommonOptions{
				Name: "karpenter-default-nodepool",
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
apiVersion: karpenter.sh/v1beta1
kind: NodePool
metadata:
  name: default
spec:
  template:
    spec:
      requirements:
        - key: kubernetes.io/arch
          operator: In
          values: ["amd64"]
        - key: kubernetes.io/os
          operator: In
          values: ["linux"]
        - key: karpenter.sh/capacity-type
          operator: In
          values: ["on-demand", "spot"]
        - key: karpenter.k8s.aws/instance-category
          operator: In
          values: ["c", "m", "r"]
        - key: karpenter.k8s.aws/instance-generation
          operator: Gt
          values: ["2"]
      nodeClassRef:
        name: default
  limits:
    cpu: 1000
  disruption:
    consolidationPolicy: WhenUnderutilized
    expireAfter: {{ .ExpireAfter }}
---
apiVersion: karpenter.k8s.aws/v1beta1
kind: EC2NodeClass
metadata:
  name: default
spec:
  amiFamily: {{ .AMIFamily }}
  role: KarpenterNodeRole-{{ .ClusterName }}
  subnetSelectorTerms:
    - tags:
        Name: eksctl-{{ .ClusterName }}-cluster/SubnetPrivate*
  securityGroupSelectorTerms:
    - tags:
        aws:eks:cluster-name: {{ .ClusterName }}
`
