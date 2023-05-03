package addon

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	eksClient *aws.EKSClient
}

func NewGetter(eksClient *aws.EKSClient) *Getter {
	return &Getter{eksClient}
}

func (g *Getter) Init() {
	if g.eksClient == nil {
		g.eksClient = aws.NewEKSClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var addon *types.Addon
	var addons []*types.Addon
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		addon, err = g.GetAddonByName(name, clusterName)
		addons = []*types.Addon{addon}
	} else {
		addons, err = g.GetAllAddons(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(addons))
}

func (g *Getter) GetAddonByName(name, clusterName string) (*types.Addon, error) {
	addon, err := g.eksClient.DescribeAddon(clusterName, name)

	return addon, aws.FormatErrorAsMessageOnly(err)
}

func (g *Getter) GetAllAddons(clusterName string) ([]*types.Addon, error) {
	addonNames, err := g.eksClient.ListAddons(clusterName)
	addons := make([]*types.Addon, 0, len(addonNames))

	if err != nil {
		return nil, err
	}

	for _, name := range addonNames {
		result, err := g.eksClient.DescribeAddon(clusterName, name)
		if err != nil {
			return nil, err
		}
		addons = append(addons, result)
	}

	return addons, nil
}
