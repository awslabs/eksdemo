package network_acl_rule

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

func (g *Getter) Get(_ string, output printer.Output, options resource.Options) error {
	naclRuleOptions, ok := options.(*NetworkAclRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NetworkAclRuleOptions")
	}

	nacls, err := g.ec2Client.DescribeNetworkAcls([]types.Filter{
		aws.NewEC2NetworkAclFilter(naclRuleOptions.NetworkAclId),
	})
	if err != nil {
		return err
	}

	if len(nacls) == 0 {
		return fmt.Errorf("network-acl %q not found", naclRuleOptions.NetworkAclId)
	}

	rules := make([]types.NetworkAclEntry, 0, len(nacls[0].Entries))
	for _, entry := range nacls[0].Entries {
		if awssdk.ToBool(entry.Egress) == naclRuleOptions.Egress {
			rules = append(rules, entry)
		}
	}

	return output.Print(os.Stdout, NewPrinter(rules))
}
