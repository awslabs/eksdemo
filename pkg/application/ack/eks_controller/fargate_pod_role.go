package eks_controller

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/cloudformation"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func fargatePodExecutionRole() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "fargate-pod-execution-role",
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

const cloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Description: Resources used by https://github.com/awslabs/karpenter
Resources:
  FargatePodExecutionRole:
    Type: "AWS::IAM::Role"
    Properties:
      RoleName: "eksdemo.{{ .ClusterName }}.fargate-pod-execution-role"
      Path: /
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
          - Effect: Allow
            Principal:
              Service:
                !Sub "eks-fargate-pods.${AWS::URLSuffix}"
            Action:
              - "sts:AssumeRole"
      ManagedPolicyArns:
        - !Sub "arn:${AWS::Partition}:iam::aws:policy/AmazonEKSFargatePodExecutionRolePolicy"
`
