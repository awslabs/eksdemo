package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type Elasticloadbalancingv2Client struct {
	*elasticloadbalancingv2.Client
}

func NewElasticloadbalancingClientv2() *Elasticloadbalancingv2Client {
	return &Elasticloadbalancingv2Client{elasticloadbalancingv2.NewFromConfig(GetConfig())}
}

func (c *Elasticloadbalancingv2Client) CreateTargetGroup(name string, port int32, protocol, targetType, vpcId string) error {
	_, err := c.Client.CreateTargetGroup(context.Background(), &elasticloadbalancingv2.CreateTargetGroupInput{
		Name:       aws.String(name),
		Port:       aws.Int32(port),
		Protocol:   types.ProtocolEnum(protocol),
		TargetType: types.TargetTypeEnum(targetType),
		VpcId:      aws.String(vpcId),
	})

	return err
}

func (c *Elasticloadbalancingv2Client) DeleteLoadBalancer(arn string) error {
	_, err := c.Client.DeleteLoadBalancer(context.Background(), &elasticloadbalancingv2.DeleteLoadBalancerInput{
		LoadBalancerArn: aws.String(arn),
	})

	return err
}

func (c *Elasticloadbalancingv2Client) DeleteTargetGroup(arn string) error {
	_, err := c.Client.DeleteTargetGroup(context.Background(), &elasticloadbalancingv2.DeleteTargetGroupInput{
		TargetGroupArn: aws.String(arn),
	})

	return err
}

func (c *Elasticloadbalancingv2Client) DescribeListeners(loadBalancerArn string) ([]types.Listener, error) {
	listeners := []types.Listener{}
	pageNum := 0

	input := elasticloadbalancingv2.DescribeListenersInput{}
	if loadBalancerArn != "" {
		input.LoadBalancerArn = aws.String(loadBalancerArn)
	}

	paginator := elasticloadbalancingv2.NewDescribeListenersPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		listeners = append(listeners, out.Listeners...)
		pageNum++
	}

	return listeners, nil
}

func (c *Elasticloadbalancingv2Client) DescribeLoadBalancers(name string) ([]types.LoadBalancer, error) {
	loadBalancers := []types.LoadBalancer{}
	pageNum := 0

	input := elasticloadbalancingv2.DescribeLoadBalancersInput{}
	if name != "" {
		input.Names = []string{name}
	}

	paginator := elasticloadbalancingv2.NewDescribeLoadBalancersPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		loadBalancers = append(loadBalancers, out.LoadBalancers...)
		pageNum++
	}

	return loadBalancers, nil
}

func (c *Elasticloadbalancingv2Client) DescribeRules(listenerArn string, ruleArns []string) ([]types.Rule, error) {
	input := elasticloadbalancingv2.DescribeRulesInput{}

	if listenerArn != "" {
		input.ListenerArn = aws.String(listenerArn)
	} else {
		input.RuleArns = ruleArns
	}

	result, err := c.Client.DescribeRules(context.Background(), &input)
	if err != nil {
		return []types.Rule{}, err
	}

	return result.Rules, nil
}

func (c *Elasticloadbalancingv2Client) DescribeTargetGroups(name, loadBalancerArn string) ([]types.TargetGroup, error) {
	targetGroups := []types.TargetGroup{}
	pageNum := 0

	input := elasticloadbalancingv2.DescribeTargetGroupsInput{}

	if name != "" {
		input.Names = []string{name}
	}

	if loadBalancerArn != "" {
		input.LoadBalancerArn = aws.String(loadBalancerArn)
	}

	paginator := elasticloadbalancingv2.NewDescribeTargetGroupsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, out.TargetGroups...)
		pageNum++
	}

	return targetGroups, nil
}

func (c *Elasticloadbalancingv2Client) DescribeTargetHealth(arn, id string) ([]types.TargetHealthDescription, error) {
	input := &elasticloadbalancingv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(arn),
	}

	if id != "" {
		input.Targets = []types.TargetDescription{
			{
				Id: aws.String(id),
			},
		}
	}

	res, err := c.Client.DescribeTargetHealth(context.Background(), input)
	if err != nil {
		return []types.TargetHealthDescription{}, err
	}

	return res.TargetHealthDescriptions, nil
}
