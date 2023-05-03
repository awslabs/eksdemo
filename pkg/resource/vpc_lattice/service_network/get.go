package service_network

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
	var serviceNetwork *vpclattice.GetServiceNetworkOutput
	var serviceNetworks []*vpclattice.GetServiceNetworkOutput
	var err error

	if id != "" {
		serviceNetwork, err = g.GetServiceNetworkById(id)
		serviceNetworks = []*vpclattice.GetServiceNetworkOutput{serviceNetwork}
	} else {
		serviceNetworks, err = g.GetAllServiceNetworks()
	}

	if err != nil {
		return aws.FormatError(err)
	}

	return output.Print(os.Stdout, NewPrinter(serviceNetworks))
}

func (g *Getter) GetAllServiceNetworks() ([]*vpclattice.GetServiceNetworkOutput, error) {
	snSummaries, err := g.vpcLatticeClient.ListServiceNetworks()
	serviceNetworks := make([]*vpclattice.GetServiceNetworkOutput, 0, len(snSummaries))

	if err != nil {
		return nil, err
	}

	for _, sn := range snSummaries {
		result, err := g.vpcLatticeClient.GetServiceNetwork(awssdk.ToString(sn.Id))
		if err != nil {
			return nil, err
		}
		serviceNetworks = append(serviceNetworks, result)
	}

	return serviceNetworks, nil
}

func (g *Getter) GetServiceNetworkById(id string) (*vpclattice.GetServiceNetworkOutput, error) {
	return g.vpcLatticeClient.GetServiceNetwork(id)
}
