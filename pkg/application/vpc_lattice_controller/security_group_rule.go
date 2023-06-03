package vpc_lattice_controller

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/awslabs/eksdemo/pkg/cloudformation"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type SecurityGroupRuleOptions struct {
	resource.CommonOptions
	VpcLatticePrefixListId *string
}

func securityGroupRule(o *GatwayApiControllerOptions) *resource.Resource {
	res := &resource.Resource{
		Options: &SecurityGroupRuleOptions{
			CommonOptions: resource.CommonOptions{
				Name: "amazon-vpc-lattice-security-group-rule",
			},
			VpcLatticePrefixListId: &o.VpcLatticePrefixListId,
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

const sqsCloudFormationTemplate = `
AWSTemplateFormatVersion: "2010-09-09"
Description: Allow traffic to Amazon EKS nodes from Amazon VPC Lattice
Resources:
  GatewayApiControllerVpcLatticeIngressRule:
    Type: AWS::EC2::SecurityGroupIngress
    Properties:
      Description: Allow traffic from Amazon VPC Lattice
      GroupId: {{ .Cluster.ResourcesVpcConfig.ClusterSecurityGroupId }}
      IpProtocol: -1
      SourcePrefixListId: {{ .VpcLatticePrefixListId }}
`
