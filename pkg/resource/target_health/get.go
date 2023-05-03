package target_health

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/target_group"
)

type Getter struct {
	elbClientv2       *aws.Elasticloadbalancingv2Client
	targetGroupGetter *target_group.Getter
}

func NewGetter(elbClientv2 *aws.Elasticloadbalancingv2Client) *Getter {
	return &Getter{elbClientv2, target_group.NewGetter(elbClientv2)}
}

func (g *Getter) Init() {
	if g.elbClientv2 == nil {
		g.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
	g.targetGroupGetter = target_group.NewGetter(g.elbClientv2)
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	thOptions, ok := options.(*TargetHealthOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to TargetHealthOptions")
	}

	targetGroup, err := g.targetGroupGetter.GetTargetGroupByName(thOptions.TargetGroupName)
	if err != nil {
		return err
	}

	targets, err := g.elbClientv2.DescribeTargetHealth(awssdk.ToString(targetGroup.TargetGroupArn), id)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(targets))
}
