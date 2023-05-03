package load_balancer

import (
	"errors"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	typesv1 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	typesv2 "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type LoadBalancers struct {
	V1 []typesv1.LoadBalancerDescription
	V2 []typesv2.LoadBalancer
}

type Getter struct {
	elbClientv1 *aws.ElasticloadbalancingClient
	elbClientv2 *aws.Elasticloadbalancingv2Client
}

func NewGetter(elbClientv1 *aws.ElasticloadbalancingClient, elbClientv2 *aws.Elasticloadbalancingv2Client) *Getter {
	return &Getter{elbClientv1, elbClientv2}
}

func (g *Getter) Init() {
	if g.elbClientv1 == nil {
		g.elbClientv1 = aws.NewElasticloadbalancingClientv1()
	}
	if g.elbClientv2 == nil {
		g.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	loadBalancers, err := g.GetLoadBalancerByName(name)
	if err != nil {
		return err
	}

	cluster := options.Common().Cluster
	if cluster != nil {
		filterByVpcId(loadBalancers, awssdk.ToString(cluster.ResourcesVpcConfig.VpcId))
	}

	return output.Print(os.Stdout, NewPrinter(loadBalancers))
}

func (g *Getter) GetAllLoadBalancers() (*LoadBalancers, error) {
	return g.GetLoadBalancerByName("")
}

func (g *Getter) GetLoadBalancerByName(name string) (*LoadBalancers, error) {
	v1, err := g.elbClientv1.DescribeLoadBalancers(name)

	// Return all errors except NotFound
	var apnfe *typesv1.AccessPointNotFoundException
	if err != nil && !errors.As(err, &apnfe) {
		return nil, err
	}

	v2, err := g.elbClientv2.DescribeLoadBalancers(name)

	// Return all errors except NotFound
	var lbnfe *typesv2.LoadBalancerNotFoundException
	if err != nil && !errors.As(err, &lbnfe) {
		return nil, err
	}

	if name != "" && len(v1) == 0 && len(v2) == 0 {
		return nil, fmt.Errorf("load balancer %q not found", name)
	}

	return &LoadBalancers{v1, v2}, nil
}

func (g *Getter) GetSecurityGroupIdsForLoadBalancer(name string) ([]string, error) {
	loadBalancers, err := g.GetLoadBalancerByName(name)
	if err != nil {
		return nil, err
	}

	// Check for the unlikely but possible scenario with elbv1 and elbv2 with same name
	if len(loadBalancers.V1) > 0 && len(loadBalancers.V2) > 0 {
		return nil, fmt.Errorf("multiple load balancers with name %q", name)
	}

	if len(loadBalancers.V2) > 0 {
		return loadBalancers.V2[0].SecurityGroups, nil
	}

	if len(loadBalancers.V1) > 0 {
		return loadBalancers.V1[0].SecurityGroups, nil
	}

	return []string{}, nil
}

func filterByVpcId(loadBalancers *LoadBalancers, id string) {
	filteredV1 := make([]typesv1.LoadBalancerDescription, 0, len(loadBalancers.V1))
	filteredV2 := make([]typesv2.LoadBalancer, 0, len(loadBalancers.V2))

	for _, v1 := range loadBalancers.V1 {
		if awssdk.ToString(v1.VPCId) == id {
			filteredV1 = append(filteredV1, v1)
		}
	}

	for _, v2 := range loadBalancers.V2 {
		if awssdk.ToString(v2.VpcId) == id {
			filteredV2 = append(filteredV2, v2)
		}
	}

	loadBalancers.V1 = filteredV1
	loadBalancers.V2 = filteredV2
}
