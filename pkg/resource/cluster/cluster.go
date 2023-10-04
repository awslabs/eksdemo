package cluster

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/nodegroup"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "cluster",
			Description: "EKS Cluster",
			Aliases:     []string{"clusters"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &eksctl.ResourceManager{
			Resource: "cluster",
			ConfigTemplate: &template.TextTemplate{
				Template: eksctl.EksctlHeader + EksctlTemplate + nodegroup.EksctlTemplate,
			},
		},
	}

	return addOptions(res)
}

const EksctlTemplate = `
addons:
- name: vpc-cni
  version: latest
  configurationValues: |-
  {{- if and (ne .KubernetesVersion "1.24") (not .DisableNetworkPolicy) }}
    enableNetworkPolicy: "true"
  {{- end }}
    env:
      ENABLE_PREFIX_DELEGATION: {{ printf "%t" .PrefixAssignment | printf "%q" }}
{{- if .IPv6 }}
- name: coredns
- name: kube-proxy
{{- end }}

cloudWatch:
  clusterLogging:
    enableTypes: ["*"]
{{- if .Fargate }}

fargateProfiles:
- name: default
  selectors:
    - namespace: fargate
{{- end }}
{{- if not .NoRoles }}

iam:
  withOIDC: true
  serviceAccounts:
{{- range .IrsaRoles }}
{{- $.IrsaTemplate.Render .Options }}
{{- end }}
{{- end }}
{{- if .IPv6 }}

kubernetesNetworkConfig:
  ipFamily: IPv6
{{- end }}
{{- if .Private }}

privateCluster:
  enabled: true
{{- end }}

vpc:
  cidr: {{ .VpcCidr }}
  hostnameType: {{ .HostnameType }}
`
