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
// Version: Latest is v0.33.2 (as of 1/31/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "autoscaling",
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

		Installer: &installer.HelmInstaller{
			ChartName:     "karpenter",
			ReleaseName:   "autoscaling-karpenter",
			RepositoryURL: "oci://public.ecr.aws/karpenter/karpenter",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
			Wait: true,
		},

		PostInstallResources: []*resource.Resource{
			karpenterDefaultNodePool(options),
		},
	}
	app.Options = options
	app.Flags = flags

	return app
}

const irsaPolicyDocument = `
Version: "2012-10-17"
Statement:
- Sid: AllowScopedEC2InstanceActions
  Effect: Allow
  Resource:
  - arn:{{ .Partition }}:ec2:{{ .Region }}::image/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}::snapshot/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:security-group/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:subnet/*
  - arn:{{ .Partition }}:ec2:{{ .Region }}:*:launch-template/*
  Action:
  - ec2:RunInstances
  - ec2:CreateFleet
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
    ForAllValues:StringEquals:
      aws:TagKeys:
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
  - ec2:DescribeAvailabilityZones
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
      iam:PassedToService: ec2.amazonaws.com
- Sid: AllowScopedInstanceProfileCreationActions
  Effect: Allow
  Resource: "*"
  Action:
  - iam:CreateInstanceProfile
  Condition:
    StringEquals:
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/topology.kubernetes.io/region: "{{ .Region }}"
    StringLike:
      aws:RequestTag/karpenter.k8s.aws/ec2nodeclass: "*"
- Sid: AllowScopedInstanceProfileTagActions
  Effect: Allow
  Resource: "*"
  Action:
  - iam:TagInstanceProfile
  Condition:
    StringEquals:
      aws:ResourceTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:ResourceTag/topology.kubernetes.io/region: "{{ .Region }}"
      aws:RequestTag/kubernetes.io/cluster/{{ .ClusterName }}: owned
      aws:RequestTag/topology.kubernetes.io/region: "{{ .Region }}"
    StringLike:
      aws:ResourceTag/karpenter.k8s.aws/ec2nodeclass: "*"
      aws:RequestTag/karpenter.k8s.aws/ec2nodeclass: "*"
- Sid: AllowScopedInstanceProfileActions
  Effect: Allow
  Resource: "*"
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
  Resource: "*"
  Action: iam:GetInstanceProfile
- Sid: AllowAPIServerEndpointDiscovery
  Effect: Allow
  Resource: arn:{{ .Partition }}:eks:{{ .Region }}:{{ .Account }}:cluster/{{ .ClusterName }}
  Action: eks:DescribeCluster
`

const valuesTemplate = `---
fullnameOverride: karpenter
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
    # -- drift is in BETA and is enabled by default.
    # Setting drift to false disables the drift disruption method to watch for drift between currently deployed nodes
    # and the desired state of nodes set in provisioners and node templates
    drift: {{ not .DisableDrift }}
`
