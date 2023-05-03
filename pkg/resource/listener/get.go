package listener

import (
	"fmt"
	"os"

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

func (g *Getter) Get(name string, output printer.Output, options resource.Options) (err error) {
	listenerOptions, ok := options.(*ListenerOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ListenerOptions")
	}

	elbs, err := g.loadBalancerGetter.GetLoadBalancerByName(listenerOptions.LoadBalancerName)
	if err != nil {
		return err
	}
	if len(elbs.V1) > 0 {
		return fmt.Errorf("%q is a classic load balancer", listenerOptions.LoadBalancerName)
	}

	lbArn := awssdk.ToString(elbs.V2[0].LoadBalancerArn)

	listeners, err := g.elbClientv2.DescribeListeners(lbArn)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(listeners))
}
