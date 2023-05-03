package target_group

import (
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	vpcLatticeClient *aws.VPCLatticeClient
}

func NewGetter(vpcLatticeClient *aws.VPCLatticeClient) *Getter {
	return &Getter{vpcLatticeClient}
}

func (g *Getter) Init() {
	if g.vpcLatticeClient == nil {
		g.vpcLatticeClient = aws.NewVPCLatticeClient()
	}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	var targetGroup *vpclattice.GetTargetGroupOutput
	var targetGroups []*vpclattice.GetTargetGroupOutput
	var err error

	if id != "" {
		targetGroup, err = g.GetTargetGroupById(id)
		targetGroups = []*vpclattice.GetTargetGroupOutput{targetGroup}
	} else {
		targetGroups, err = g.GetAllTargetGroups()
	}

	if err != nil {
		return aws.FormatError(err)
	}

	return output.Print(os.Stdout, NewPrinter(targetGroups))
}

func (g *Getter) GetAllTargetGroups() ([]*vpclattice.GetTargetGroupOutput, error) {
	tgSummaries, err := g.vpcLatticeClient.ListTargetGroups()
	targetGroups := make([]*vpclattice.GetTargetGroupOutput, 0, len(tgSummaries))

	if err != nil {
		return nil, err
	}

	for _, sn := range tgSummaries {
		result, err := g.vpcLatticeClient.GetTargetGroup(awssdk.ToString(sn.Id))
		if err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, result)
	}

	return targetGroups, nil
}

func (g *Getter) GetTargetGroupById(id string) (*vpclattice.GetTargetGroupOutput, error) {
	return g.vpcLatticeClient.GetTargetGroup(id)
}
