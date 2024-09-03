package userprofile

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	sagemakerClient *aws.SageMakerClient
}

func NewGetter(sagemakerClient *aws.SageMakerClient) *Getter {
	return &Getter{sagemakerClient}
}

func (g *Getter) Init() {
	if g.sagemakerClient == nil {
		g.sagemakerClient = aws.NewSageMakerClient()
	}
}

func (g *Getter) Get(userProfileName string, output printer.Output, o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to domain.Options")
	}

	var userProfile *sagemaker.DescribeUserProfileOutput
	var userProfiles []*sagemaker.DescribeUserProfileOutput
	var err error

	switch {
	case userProfileName != "":
		userProfile, err = g.GetUserProfileByName(userProfileName)
		userProfiles = []*sagemaker.DescribeUserProfileOutput{userProfile}
	case options.DomainID != "":
		userProfiles, err = g.GetUserProfilesByDomainID(options.DomainID)
	default:
		userProfiles, err = g.GetAllUserProfiles()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(userProfiles))
}

func (g *Getter) GetAllUserProfiles() ([]*sagemaker.DescribeUserProfileOutput, error) {
	return g.getUserProfiles("", "")
}

func (g *Getter) GetUserProfilesByDomainID(domainID string) ([]*sagemaker.DescribeUserProfileOutput, error) {
	return g.getUserProfiles(domainID, "")

}

func (g *Getter) GetUserProfileByName(userProfileName string) (*sagemaker.DescribeUserProfileOutput, error) {
	profileDetails, err := g.sagemakerClient.ListUserProfiles("", "")
	if err != nil {
		return nil, err
	}

	found := []types.UserProfileDetails{}

	for _, up := range profileDetails {
		if strings.EqualFold(userProfileName, awssdk.ToString(up.UserProfileName)) {
			found = append(found, up)
		}
	}

	if len(found) == 0 {
		return nil, &resource.NotFoundByNameError{Type: "sagemaker-user-profile", Name: userProfileName}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple sagemaker user profiles found with name: %s", userProfileName)
	}

	userProfile, err := g.sagemakerClient.DescribeUserProfile(awssdk.ToString(found[0].DomainId), userProfileName)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (g *Getter) getUserProfiles(domainID, userProfileNameContains string) ([]*sagemaker.DescribeUserProfileOutput, error) {
	profileDetails, err := g.sagemakerClient.ListUserProfiles(domainID, userProfileNameContains)
	if err != nil {
		return nil, err
	}

	userprofiles := make([]*sagemaker.DescribeUserProfileOutput, 0, len(profileDetails))

	for _, up := range profileDetails {
		result, err := g.sagemakerClient.DescribeUserProfile(
			awssdk.ToString(up.DomainId),
			awssdk.ToString(up.UserProfileName),
		)
		if err != nil {
			return nil, err
		}
		userprofiles = append(userprofiles, result)
	}

	return userprofiles, nil
}
