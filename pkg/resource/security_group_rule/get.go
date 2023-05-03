package security_group_rule

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
	"github.com/awslabs/eksdemo/pkg/resource/network_interface"
)

type Getter struct {
	ec2Client *aws.EC2Client
	eniGetter *network_interface.Getter
	elbGetter *load_balancer.Getter
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client, network_interface.NewGetter(ec2Client),
		load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
	g.eniGetter = network_interface.NewGetter(g.ec2Client)
	g.elbGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sgrOptions, ok := options.(*SecurityGroupRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SecurityGroupRuleOptions")
	}

	var err error
	var securityGroupRules []types.SecurityGroupRule

	if sgrOptions.SecurityGroupId != "" {
		securityGroupRules, err = g.GetSecurityGroupRulesBySecurityGroupId(sgrOptions.SecurityGroupId)
	} else if sgrOptions.LoadBalancerName != "" {
		securityGroupRules, err = g.GetSecurityGroupRulesByLoadBalancerName(sgrOptions.LoadBalancerName)
	} else if sgrOptions.NetworkInterfaceId != "" {
		securityGroupRules, err = g.GetSecurityGroupRulesByNetworkInterfaceId(sgrOptions.NetworkInterfaceId)
	} else {
		securityGroupRules, err = g.GetSecurityGroupRulesById(id)
	}

	if err != nil {
		return err
	}

	if sgrOptions.Ingress {
		securityGroupRules = filterResults(securityGroupRules, false)
	} else if sgrOptions.Egress {
		securityGroupRules = filterResults(securityGroupRules, true)
	}

	return output.Print(os.Stdout, NewPrinter(securityGroupRules, options.Common().ClusterName))
}

func (g *Getter) GetSecurityGroupRulesById(id string) ([]types.SecurityGroupRule, error) {
	return g.ec2Client.DescribeSecurityGroupRules(id, "")
}

func (g *Getter) GetSecurityGroupRulesByLoadBalancerName(name string) ([]types.SecurityGroupRule, error) {
	sgIds, err := g.elbGetter.GetSecurityGroupIdsForLoadBalancer(name)
	if err != nil {
		return nil, err
	}

	eniRules := []types.SecurityGroupRule{}

	for _, id := range sgIds {
		sgr, err := g.ec2Client.DescribeSecurityGroupRules("", id)
		if err != nil {
			return nil, err
		}
		eniRules = append(eniRules, sgr...)
	}

	return eniRules, nil
}

func (g *Getter) GetSecurityGroupRulesByNetworkInterfaceId(eniId string) ([]types.SecurityGroupRule, error) {
	networkInterface, err := g.eniGetter.GetNetworkInterfaceById(eniId)
	if err != nil {
		return nil, err
	}
	eniRules := []types.SecurityGroupRule{}

	for _, groupIdentifier := range networkInterface.Groups {
		sgr, err := g.ec2Client.DescribeSecurityGroupRules("", awssdk.ToString(groupIdentifier.GroupId))
		if err != nil {
			return nil, err
		}
		eniRules = append(eniRules, sgr...)
	}

	return eniRules, nil
}

func (g *Getter) GetSecurityGroupRulesBySecurityGroupId(securityGroupId string) ([]types.SecurityGroupRule, error) {
	return g.ec2Client.DescribeSecurityGroupRules("", securityGroupId)
}

func filterResults(rules []types.SecurityGroupRule, egress bool) []types.SecurityGroupRule {
	filtered := make([]types.SecurityGroupRule, 0, len(rules))

	if egress {
		for _, rule := range rules {
			if awssdk.ToBool(rule.IsEgress) {
				filtered = append(filtered, rule)
			}
		}
	} else {
		for _, rule := range rules {
			if !awssdk.ToBool(rule.IsEgress) {
				filtered = append(filtered, rule)
			}
		}
	}

	return filtered
}
