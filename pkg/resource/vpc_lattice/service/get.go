package service

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
	var service *vpclattice.GetServiceOutput
	var services []*vpclattice.GetServiceOutput
	var err error

	if id != "" {
		service, err = g.GetServiceById(id)
		services = []*vpclattice.GetServiceOutput{service}
	} else {
		services, err = g.GetAllServices()
	}

	if err != nil {
		return aws.FormatError(err)
	}

	return output.Print(os.Stdout, NewPrinter(services))
}

func (g *Getter) GetAllServices() ([]*vpclattice.GetServiceOutput, error) {
	summaries, err := g.vpcLatticeClient.ListServices()
	services := make([]*vpclattice.GetServiceOutput, 0, len(summaries))

	if err != nil {
		return nil, err
	}

	for _, s := range summaries {
		result, err := g.vpcLatticeClient.GetService(awssdk.ToString(s.Id))
		if err != nil {
			return nil, err
		}
		services = append(services, result)
	}

	return services, nil
}

func (g *Getter) GetServiceById(id string) (*vpclattice.GetServiceOutput, error) {
	return g.vpcLatticeClient.GetService(id)
}
