package nodegroup

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "nodegroup",
			Description: "EKS Managed Nodegroup",
			Aliases:     []string{"nodegroups", "mng", "ng"},
			CreateArgs:  []string{"NAME"},
			Args:        []string{"NAME"},
		},

		Getter: &Getter{},

		Manager: &Manager{
			Eksctl: &eksctl.ResourceManager{
				Resource: "nodegroup",
				ConfigTemplate: &template.TextTemplate{
					Template: eksctl.EksctlHeader + EksctlTemplate,
				},
				ApproveDelete: true,
				CreateFlags:   []string{"--install-nvidia-plugin=false", "--install-neuron-plugin=false"},
			},
		},
	}

	res.Options, res.CreateFlags, res.UpdateFlags = NewOptions()

	return res
}

const EksctlTemplate = `
managedNodeGroups:
- name: {{ .NodegroupName }}
{{- if .AMI }}
  ami: {{ .AMI }}
{{- end }}
  amiFamily: {{ .OperatingSystem }}
  desiredCapacity: {{ .DesiredCapacity }}
  iam:
    attachPolicyARNs:
    - arn:{{ .Partition }}:iam::aws:policy/AmazonEKSWorkerNodePolicy
    - arn:{{ .Partition }}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly
    - arn:{{ .Partition }}:iam::aws:policy/AmazonSSMManagedInstanceCore
{{- if .Spot }}
  instanceSelector:
    vCPUs: {{ .SpotvCPUs }}
    memory: {{ .SpotMemory | toString | printf "%q" }}
{{- else }}
  instanceType: {{ .InstanceType }}
{{- end }}
  minSize: {{ .MinSize }}
  maxSize: {{ .MaxSize }}
  volumeSize: {{ .VolumeSize }}
  volumeType: {{ .VolumeType }}
{{- if eq .VolumeType "io2" }}
  volumeIOPS: 16000
{{- end }}
{{- if and .AMI (ne .OperatingSystem "AmazonLinux2023") }}
  overrideBootstrapCommand: |
    #!/bin/bash
    /etc/eks/bootstrap.sh {{ .ClusterName }}
{{- end }}
  privateNetworking: true
  spot: {{ .Spot }}
{{- range .Taints }}
  taints:
  - key: {{ .Key }}
    value: {{ .Value | printf "%q" }}
    effect: {{ .Effect }}
{{- end }}
`
