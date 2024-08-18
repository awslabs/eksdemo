package karpenter

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/cloudformation"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func karpenterSqsQueue() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "karpenter-sqs-queue",
		},

		Manager: &cloudformation.ResourceManager{
			Capabilities: []types.Capability{types.CapabilityCapabilityNamedIam},
			Template: &template.TextTemplate{
				Template: sqsCloudFormationTemplate,
			},
		},
	}
	return res
}

// https://github.com/aws/karpenter/blob/main/website/content/en/preview/getting-started/getting-started-with-karpenter/cloudformation.yaml#L80
const sqsCloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Description: Resources used by https://github.com/aws/karpenter
Resources:
  KarpenterInterruptionQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: karpenter-{{ .ClusterName }}
      MessageRetentionPeriod: 300
      SqsManagedSseEnabled: true
  KarpenterInterruptionQueuePolicy:
    Type: AWS::SQS::QueuePolicy
    Properties:
      Queues:
        - !Ref KarpenterInterruptionQueue
      PolicyDocument:
        Id: EC2InterruptionPolicy
        Statement:
          - Effect: Allow
            Principal:
              Service:
                - events.amazonaws.com
                - sqs.amazonaws.com
            Action: sqs:SendMessage
            Resource: !GetAtt KarpenterInterruptionQueue.Arn
          - Sid: DenyHTTP
            Effect: Deny
            Action: sqs:*
            Resource: !GetAtt KarpenterInterruptionQueue.Arn
            Condition:
              Bool:
                aws:SecureTransport: false
            Principal: "*"
  ScheduledChangeRule:
    Type: 'AWS::Events::Rule'
    Properties:
      Name: {{ printf "karpenter-%s-ScheduledChange" .ClusterName | .TruncateUnique 64 }}
      EventPattern:
        source:
          - aws.health
        detail-type:
          - AWS Health Event
      Targets:
        - Id: KarpenterInterruptionQueueTarget
          Arn: !GetAtt KarpenterInterruptionQueue.Arn
  SpotInterruptionRule:
    Type: 'AWS::Events::Rule'
    Properties:
      Name: {{ printf "karpenter-%s-SpotInterruption" .ClusterName | .TruncateUnique 64 }}
      EventPattern:
        source:
          - aws.ec2
        detail-type:
          - EC2 Spot Instance Interruption Warning
      Targets:
        - Id: KarpenterInterruptionQueueTarget
          Arn: !GetAtt KarpenterInterruptionQueue.Arn
  RebalanceRule:
    Type: 'AWS::Events::Rule'
    Properties:
      Name: {{ printf "karpenter-%s-Rebalance" .ClusterName | .TruncateUnique 64 }}
      EventPattern:
        source:
          - aws.ec2
        detail-type:
          - EC2 Instance Rebalance Recommendation
      Targets:
        - Id: KarpenterInterruptionQueueTarget
          Arn: !GetAtt KarpenterInterruptionQueue.Arn
  InstanceStateChangeRule:
    Type: 'AWS::Events::Rule'
    Properties:
      Name: {{ printf "karpenter-%s-InstanceStateChange" .ClusterName | .TruncateUnique 64 }}
      EventPattern:
        source:
          - aws.ec2
        detail-type:
          - EC2 Instance State-change Notification
      Targets:
        - Id: KarpenterInterruptionQueueTarget
          Arn: !GetAtt KarpenterInterruptionQueue.Arn
`
