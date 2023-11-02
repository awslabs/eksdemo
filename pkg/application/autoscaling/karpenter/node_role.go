package karpenter

import (
	"github.com/awslabs/eksdemo/pkg/cloudformation"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func karpenterNodeRole() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "karpenter-node-role",
		},

		Manager: &cloudformation.ResourceManager{
			Capabilities: []types.Capability{types.CapabilityCapabilityNamedIam},
			Template: &template.TextTemplate{
				Template: cloudFormationTemplate,
			},
		},
	}
	return res
}

// https://github.com/aws/karpenter/blob/main/website/content/en/preview/getting-started/getting-started-with-karpenter/cloudformation.yaml
const cloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Description: Resources used by https://github.com/aws/karpenter
Parameters:
  ClusterName:
    Type: String
    Description: "EKS cluster name"
    Default: "{{ .ClusterName }}"
Resources:
  KarpenterNodeRole:
    Type: "AWS::IAM::Role"
    Properties:
      RoleName: !Sub "KarpenterNodeRole-${ClusterName}"
      Path: /
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service:
                !Sub "ec2.${AWS::URLSuffix}"
            Action:
              - "sts:AssumeRole"
      ManagedPolicyArns:
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonSSMManagedInstanceCore"
`
