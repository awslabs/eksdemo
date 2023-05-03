package vpc_summary

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
)

type Getter struct {
	ec2Client *aws.EC2Client
	lbGetter  *load_balancer.Getter
}

type VpcSummary struct {
	eips              []types.Address
	endpoints         []types.VpcEndpoint
	instances         []types.Reservation
	internetGWs       []types.InternetGateway
	loadBalancers     *load_balancer.LoadBalancers
	natGWs            []types.NatGateway
	networkInterfaces []types.NetworkInterface
	routeTables       []types.RouteTable
	securityGroups    []types.SecurityGroup
	subnets           []types.Subnet
	vpcs              []types.Vpc
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{
		ec2Client,
		load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2()),
	}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
	g.lbGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	vpcSummaryOptions, ok := options.(*VpcSummaryOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to VpcSummaryOptions")
	}

	cluster := options.Common().Cluster
	igwfilter := []types.Filter{}
	filters := []types.Filter{}
	vpcFilter := false

	if id != "" {
		igwfilter = append(filters, aws.NewEC2InternetGatewayVpcFilter(id))
		filters = append(filters, aws.NewEC2VpcFilter(id))
	}

	if cluster != nil {
		igwfilter = append(filters, aws.NewEC2InternetGatewayVpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	eips, err := g.ec2Client.DescribeAddresses(nil)
	if err != nil {
		return err
	}

	endpoints, err := g.ec2Client.DescribeVpcEndpoints(filters)
	if err != nil {
		return err
	}

	instances, err := g.ec2Client.DescribeInstances(append(filters, aws.NewEC2InstanceStateFilter([]string{"running"})))
	if err != nil {
		return err
	}

	internetGWs, err := g.ec2Client.DescribeInternetGateways(igwfilter)
	if err != nil {
		return err
	}

	loadBalancers, err := g.lbGetter.GetAllLoadBalancers()
	if err != nil {
		return err
	}

	natGWs, err := g.ec2Client.DescribeNATGateways(filters)
	if err != nil {
		return err
	}

	networkInterfaces, err := g.ec2Client.DescribeNetworkInterfaces(filters)
	if err != nil {
		return err
	}

	routeTables, err := g.ec2Client.DescribeRouteTables(filters)
	if err != nil {
		return err
	}

	securityGroups, err := g.ec2Client.DescribeSecurityGroups(filters, nil)
	if err != nil {
		return err
	}

	subnets, err := g.ec2Client.DescribeSubnets(filters)
	if err != nil {
		return err
	}

	vpcs, err := g.ec2Client.DescribeVpcs(filters)
	if err != nil {
		return err
	}

	summary := &VpcSummary{
		eips,
		endpoints,
		instances,
		internetGWs,
		loadBalancers,
		natGWs,
		networkInterfaces,
		routeTables,
		securityGroups,
		subnets,
		vpcs,
	}

	if len(filters) > 0 {
		vpcFilter = true
	}

	return output.Print(os.Stdout, NewPrinter(summary, vpcFilter, vpcSummaryOptions.ShowIds))
}
