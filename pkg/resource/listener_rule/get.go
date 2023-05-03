package listener_rule

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
)

type Getter struct {
	elbClientv2        *aws.Elasticloadbalancingv2Client
	loadBalancerGetter *load_balancer.Getter
}

func NewGetter(elbClientv2 *aws.Elasticloadbalancingv2Client) *Getter {
	return &Getter{elbClientv2, load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())}
}

func (g *Getter) Init() {
	if g.elbClientv2 == nil {
		g.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
	g.loadBalancerGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) (err error) {
	lrOptions, ok := options.(*ListenerRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ListenerRuleOptions")
	}

	elbs, err := g.loadBalancerGetter.GetLoadBalancerByName(lrOptions.LoadBalancerName)
	if err != nil {
		return err
	}
	if len(elbs.V1) > 0 {
		return fmt.Errorf("%q is a classic load balancer", lrOptions.LoadBalancerName)
	}

	lbArn := awssdk.ToString(elbs.V2[0].LoadBalancerArn)
	listernArn := ""
	ruleArn := ""

	if id == "" {
		listernArn = strings.Replace(lbArn, ":loadbalancer/", ":listener/", 1) + "/" + lrOptions.ListenerId
	} else {
		ruleArn = strings.Replace(lbArn, ":loadbalancer/", ":listener-rule/", 1) + "/" + lrOptions.ListenerId + "/" + id
	}

	rules, err := g.elbClientv2.DescribeRules(listernArn, []string{ruleArn})
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(rules))
}
