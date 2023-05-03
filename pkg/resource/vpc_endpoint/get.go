package vpc_endpoint

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	ec2Client *aws.EC2Client
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	vpcOptions, ok := options.(*VpcEndpointOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to VpcEndpointOptions")
	}

	cluster := options.Common().Cluster
	filters := []types.Filter{}

	if id != "" {
		filters = append(filters, aws.NewEC2VpcEndpointFilter(id))
	}

	if cluster != nil {
		filters = append(filters, aws.NewEC2VpcFilter(awssdk.ToString(cluster.ResourcesVpcConfig.VpcId)))
	}

	if vpcOptions.VpcId != "" {
		filters = append(filters, aws.NewEC2VpcFilter(vpcOptions.VpcId))
	}

	vpcEndpoints, err := g.ec2Client.DescribeVpcEndpoints(filters)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(vpcEndpoints))
}
