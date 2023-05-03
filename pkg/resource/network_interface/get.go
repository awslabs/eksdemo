package network_interface

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/load_balancer"
)

type Getter struct {
	ec2Client          *aws.EC2Client
	loadBalancerGetter *load_balancer.Getter
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client, load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
	g.loadBalancerGetter = load_balancer.NewGetter(aws.NewElasticloadbalancingClientv1(), aws.NewElasticloadbalancingClientv2())
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	eniOptions, ok := options.(*NetworkInterfaceOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NetworkInterfaceOptions")
	}

	var eni types.NetworkInterface
	var enis []types.NetworkInterface
	var err error

	if id != "" {
		eni, err = g.GetNetworkInterfaceById(id)
		enis = []types.NetworkInterface{eni}
	} else {
		enis, err = g.GetNetworkInterfacesWithOptions(eniOptions)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(enis, aws.AccountId()))
}

func (g *Getter) GetNetworkInterfaceById(id string) (types.NetworkInterface, error) {
	filters := []types.Filter{aws.NewEC2NetworkInterfaceFilter(id)}

	networkInterfaces, err := g.ec2Client.DescribeNetworkInterfaces(filters)
	if err != nil {
		return types.NetworkInterface{}, err
	}

	if len(networkInterfaces) == 0 {
		return types.NetworkInterface{}, resource.NotFoundError(fmt.Sprintf("network-interface %q not found", id))
	}

	return networkInterfaces[0], nil
}

func (g *Getter) GetNetworkInterfacesWithOptions(eniOptions *NetworkInterfaceOptions) ([]types.NetworkInterface, error) {
	cluster := eniOptions.Cluster
	filters := []types.Filter{}

	if cluster != nil {
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	if eniOptions.InstanceId != "" {
		filters = append(filters, aws.NewEC2NetworkInterfaceInstanceFilter(eniOptions.InstanceId))

	}

	if eniOptions.IpAddress != "" {
		filters = append(filters, aws.NewEC2NetworkInterfaceIpFilter(eniOptions.IpAddress))

	}

	if eniOptions.LoadBalancerName != "" {
		elbs, err := g.loadBalancerGetter.GetLoadBalancerByName(eniOptions.LoadBalancerName)
		if err != nil {
			return nil, err
		}

		var description string
		// Identify the ENIs for a LoadBalancer using Description as described below
		// https://aws.amazon.com/premiumsupport/knowledge-center/elb-find-load-balancer-IP/
		if len(elbs.V1) > 0 {
			description = "ELB " + *elbs.V1[0].LoadBalancerName
		} else if len(elbs.V2) > 0 {
			elb := elbs.V2[0]
			lbName := awssdk.ToString(elb.LoadBalancerName)
			lbArn := awssdk.ToString(elb.LoadBalancerArn)
			lbId := lbArn[strings.LastIndex(lbArn, "/")+1:]

			switch string(elb.Type) {
			case "application":
				description = fmt.Sprintf("ELB app/%s/%s", lbName, lbId)
			case "network":
				description = fmt.Sprintf("ELB net/%s/%s", lbName, lbId)
			default:
				return nil, fmt.Errorf("load balancer type %q not supported", string(elb.Type))
			}
		}
		filters = append(filters, aws.NewEC2DescriptionFilter(description))
	}

	if eniOptions.SecurityGroupId != "" {
		filters = append(filters, aws.NewEC2SecurityGroupFilter(eniOptions.SecurityGroupId))
	}

	return g.ec2Client.DescribeNetworkInterfaces(filters)
}
