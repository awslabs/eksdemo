package fargate_profile

import (
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
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
	var profile *types.FargateProfile
	var profiles []*types.FargateProfile
	var err error

	clusterName := options.Common().ClusterName

	if name != "" {
		profile, err = g.GetProfileByName(name, clusterName)
		profiles = []*types.FargateProfile{profile}
	} else {
		profiles, err = g.GetAllProfiles(clusterName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(profiles))
}

func (g *Getter) GetAllProfiles(clusterName string) ([]*types.FargateProfile, error) {
	profileNames, err := g.eksClient.ListFargateProfiles(clusterName)
	profiles := make([]*types.FargateProfile, 0, len(profileNames))

	if err != nil {
		return nil, err
	}

	for _, name := range profileNames {
		result, err := g.eksClient.DescribeFargateProfile(clusterName, name)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, result)
	}

	return profiles, nil
}

func (g *Getter) GetProfileByName(name, clusterName string) (*types.FargateProfile, error) {
	profile, err := g.eksClient.DescribeFargateProfile(clusterName, name)

	return profile, aws.FormatErrorAsMessageOnly(err)
}
