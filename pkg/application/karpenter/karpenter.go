package karpenter

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/iam_auth"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/resource/service_linked_role"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://karpenter.sh/docs/
// GitHub:  https://github.com/awslabs/karpenter
// Helm:    https://github.com/awslabs/karpenter/tree/main/charts/karpenter
// Repo:    https://gallery.ecr.aws/karpenter/controller
// Version: Latest is v1.5.0 (as of 6/7/25)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Name:        "karpenter",
			Description: "Karpenter Node Autoscaling",
		},

		Dependencies: []*resource.Resource{
			service_linked_role.NewResourceWithOptions(&service_linked_role.ServiceLinkedRoleOptions{
				CommonOptions: resource.CommonOptions{
					Name: "ec2-spot-service-linked-role",
				},
				RoleName:    "AWSServiceRoleForEC2Spot",
				ServiceName: "spot.amazonaws.com",
			}),
			irsa.NewResourceWithOptions(&irsa.IrsaOptions{
				CommonOptions: resource.CommonOptions{
					Name: "karpenter-irsa",
				},
				PolicyType: irsa.PolicyDocument,
				PolicyDocTemplate: &template.TextTemplate{
					Template: irsaPolicyDocument,
				},
			}),
			karpenterNodeRole(),
			karpenterSqsQueue(),
			iam_auth.NewResourceWithOptions(&iam_auth.IamAuthOptions{
				CommonOptions: resource.CommonOptions{
					Name: "karpenter-node-iam-auth",
				},
				Groups: []string{"system:bootstrappers", "system:nodes"},
				RoleName: &template.TextTemplate{
					Template: "KarpenterNodeRole-{{ .ClusterName }}",
				},
				Username: "system:node:{{EC2PrivateDNSName}}",
			}),
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "karpenter",
			ReleaseName:   "karpenter",
			RepositoryURL: "oci://public.ecr.aws/karpenter/karpenter",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			Wait: true,
		},

		Options: options,

		PostInstallResources: []*resource.Resource{
			karpenterDefaultNodePool(options),
		},
	}
}

const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Sid: AllowScopedEC2InstanceAccessActions
  Effect: Allow
  Resource:
  - arn:{{ .Partition }}:ec2:{{ .Region }}::image/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}::snapshot/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:security-group/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:subnet/*
  Action:
  - ec2:RunInstances
  - ec2:CreateFleet
- Sid: AllowScopedEC2LaunchTemplateAccessActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:ec2:{{ .Region }}:*:launch-template/*
  Action:
  - ec2:RunInstances
  - ec2:CreateFleet
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
    StringLike:
      aws:ResourceTag/karpenter.sh/nodepool: "*"
- Sid: AllowScopedEC2InstanceActionsWithTags
  Effect: Allow
  Resource:
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:fleet/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:instance/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:volume/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:network-interface/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:launch-template/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:spot-instances-request/*
  Action:
  - ec2:RunInstances
  - ec2:CreateFleet
  - ec2:CreateLaunchTemplate
  Condition:
    StringEquals:
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/eks:eks-cluster-name: {{ .ClusterName }}
    StringLike:
      aws:RequestTag/karpenter.sh/nodepool: "*"
- Sid: AllowScopedResourceCreationTagging
  Effect: Allow
  Resource:
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:fleet/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:instance/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:volume/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:network-interface/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:launch-template/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:spot-instances-request/*
  Action: ec2:CreateTags
  Condition:
    StringEquals:
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/eks:eks-cluster-name: {{ .ClusterName }}
      ec2:CreateAction:
      - RunInstances
      - CreateFleet
      - CreateLaunchTemplate
    StringLike:
      aws:RequestTag/karpenter.sh/nodepool: "*"
- Sid: AllowScopedResourceTagging
  Effect: Allow
  Resource: arn:{{ .Partition }}:ec2:{{ .Region }}:*:instance/*
  Action: ec2:CreateTags
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
    StringLike:
      aws:ResourceTag/karpenter.sh/nodepool: "*"
    StringEqualsIfExists:
      aws:RequestTag/eks:eks-cluster-name: {{ .ClusterName }}
    ForAllValues:StringEquals:
      aws:TagKeys:
      - eks:eks-cluster-name
      - karpenter.sh/nodeclaim
      - Name
- Sid: AllowScopedDeletion
  Effect: Allow
  Resource:
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:instance/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:launch-template/*
  Action:
  - ec2:TerminateInstances
  - ec2:DeleteLaunchTemplate
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
    StringLike:
      aws:ResourceTag/karpenter.sh/nodepool: "*"
- Sid: AllowRegionalReadActions
  Effect: Allow
  Resource: "*"
  Action:
  - ec2:DescribeImages
  - ec2:DescribeInstances
  - ec2:DescribeInstanceTypeOfferings
  - ec2:DescribeInstanceTypes
  - ec2:DescribeLaunchTemplates
  - ec2:DescribeSecurityGroups
  - ec2:DescribeSpotPriceHistory
  - ec2:DescribeSubnets
  Condition:
    StringEquals:
      aws:RequestedRegion: "{{ .Region }}"
- Sid: AllowSSMReadActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:ssm:{{ .Region }}::parameter/aws/service/*
  Action:
  - ssm:GetParameter
- Sid: AllowPricingReadActions
  Effect: Allow
  Resource: "*"
  Action:
  - pricing:GetProducts
- Sid: AllowInterruptionQueueActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:sqs:{{ .Region }}:{{ .Account }}:karpenter-{{ .ClusterName }}
  Action:
  - sqs:DeleteMessage
  - sqs:GetQueueUrl
  - sqs:ReceiveMessage
- Sid: AllowPassingInstanceRole
  Effect: Allow
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:role/KarpenterNodeRole-{{ .ClusterName }}
  Action: iam:PassRole
  Condition:
    StringEquals:
      iam:PassedToService:
      - ec2.amazonaws.com
      - ec2.amazonaws.com.cn
- Sid: AllowScopedInstanceProfileCreationActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:instance-profile/*
  Action:
  - iam:CreateInstanceProfile
  Condition:
    StringEquals:
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/eks:eks-cluster-name: {{ .ClusterName }}
      aws:RequestTag/topology.kubernetes.io/region: "{{ .Region }}"
    StringLike:
      aws:RequestTag/karpenter.k8s.aws/ec2nodeclass: "*"
- Sid: AllowScopedInstanceProfileTagActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:instance-profile/*
  Action:
  - iam:TagInstanceProfile
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:ResourceTag/topology.kubernetes.io/region: "{{ .Region }}"
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/eks:eks-cluster-name: {{ .ClusterName }}
      aws:RequestTag/topology.kubernetes.io/region: "{{ .Region }}"
    StringLike:
      aws:ResourceTag/karpenter.k8s.aws/ec2nodeclass: "*"
      aws:RequestTag/karpenter.k8s.aws/ec2nodeclass: "*"
- Sid: AllowScopedInstanceProfileActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:instance-profile/*
  Action:
  - iam:AddRoleToInstanceProfile
  - iam:RemoveRoleFromInstanceProfile
  - iam:DeleteInstanceProfile
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:ResourceTag/topology.kubernetes.io/region: "{{ .Region }}"
    StringLike:
      aws:ResourceTag/karpenter.k8s.aws/ec2nodeclass: "*"
- Sid: AllowInstanceProfileReadActions
  Effect: Allow
  Resource: arn:{{ .Partition }}:iam::{{ .Account }}:instance-profile/*
  Action: iam:GetInstanceProfile
- Sid: AllowAPIServerEndpointDiscovery
  Effect: Allow
  Resource: arn:{{ .Partition }}:eks:{{ .Region }}:{{ .Account }}:cluster/{{ .ClusterName }}
  Action: eks:DescribeCluster
`

// https://github.com/aws/karpenter-provider-aws/blob/main/charts/karpenter/values.yaml
const valuesTemplate = `---
serviceAccount:
  name: {{ .ServiceAccount }}
  annotations:
    {{ .IrsaAnnotation }}
replicas: {{ .Replicas }}
controller:
  image:
    tag: {{ .Version }}
  resources:
    requests:
      cpu: "1"
      memory: "1Gi"
settings:
  clusterName: {{ .ClusterName }}
  interruptionQueue: karpenter-{{ .ClusterName }}
  featureGates:
    # -- spotToSpotConsolidation is ALPHA and is disabled by default.
    # Setting this to true will enable spot replacement consolidation for both single and multi-node consolidation.
    spotToSpotConsolidation: {{ .EnableSpotToSpot }}
`
