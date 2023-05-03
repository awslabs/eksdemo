package crossplane

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cloudformation"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/awslabs/eksdemo/pkg/template"
)

func crossplaneIrsa(options application.Options) *resource.Resource {
	o := options.Common()

	return &resource.Resource{
		Options: &irsa.IrsaOptions{
			CommonOptions: resource.CommonOptions{
				ClusterName:    o.ClusterName,
				Name:           fmt.Sprintf("crossplane-%s-irsa", o.Namespace),
				Namespace:      o.Namespace,
				ServiceAccount: o.ServiceAccount,
			},
		},

		Manager: &cloudformation.ResourceManager{
			Capabilities: []types.Capability{types.CapabilityCapabilityNamedIam},
			Template: &template.TextTemplate{
				Template: cloudFormationTemplate,
			},
		},
	}
}

const cloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Resources:
  CrossplaneIRSA:
    Type: "AWS::IAM::Role"
    Properties:
      RoleName: {{ .RoleName }}
      Path: /
      AssumeRolePolicyDocument:
        Version: "2012-10-17"
        Statement:
        - Effect: Allow
          Principal:
            Federated:
              !Sub "arn:${AWS::Partition}:iam::{{ .Account }}:oidc-provider/{{ .ClusterOIDCProvider }}"
          Action:
          - sts:AssumeRoleWithWebIdentity
          Condition:
            StringLike:
              "{{ .ClusterOIDCProvider }}:sub": "system:serviceaccount:{{ .Namespace }}:provider-aws-*"
      ManagedPolicyArns:
      - !Sub "arn:${AWS::Partition}:iam::aws:policy/AdministratorAccess"
`
