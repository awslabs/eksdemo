package addon

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewVersionsResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "addon-versions",
			Description: "EKS Managed Addon Versions",
			Args:        []string{"NAME"},
		},

		Getter: &VersionGetter{},
	}

	res.Options = &resource.CommonOptions{}
	res.CreateFlags = cmd.Flags{}

	return res
}

type VersionGetter struct {
	eksClient *aws.EKSClient
}

func NewVersionGetter(eksClient *aws.EKSClient) *Getter {
	return &Getter{eksClient}
}

func (g *VersionGetter) Init() {
	if g.eksClient == nil {
		g.eksClient = aws.NewEKSClient()
	}
}

func (g *VersionGetter) Get(name string, output printer.Output, options resource.Options) error {
	var addonVersions []types.AddonInfo
	var err error

	k8sversion := options.Common().KubernetesVersion

	if name != "" {
		addonVersions, err = g.GetAddonVersionsByName(name, k8sversion)
	} else {
		addonVersions, err = g.GetAllAddonVersions(k8sversion)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewVersionPrinter(addonVersions))
}

func (g *VersionGetter) GetAddonVersionsByName(name, k8sversion string) ([]types.AddonInfo, error) {
	addonVersions, err := g.eksClient.DescribeAddonVersions(name, k8sversion)
	if err != nil {
		return nil, err
	}

	return addonVersions, nil
}

func (g *VersionGetter) GetAllAddonVersions(k8sversion string) ([]types.AddonInfo, error) {
	addonVersions, err := g.eksClient.DescribeAddonVersions("", k8sversion)
	if err != nil {
		return nil, err
	}

	return addonVersions, nil
}
