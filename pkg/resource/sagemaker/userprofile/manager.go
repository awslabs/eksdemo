package userprofile

import (
	"errors"
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Manager struct {
	resource.CreateNotSupported
	resource.DryRun
	resource.UpdateNotSupported
	sagemakerClient   *aws.SageMakerClient
	userProfileGetter *Getter
}

func (m *Manager) Init() {
	if m.sagemakerClient == nil {
		m.sagemakerClient = aws.NewSageMakerClient()
	}
	m.userProfileGetter = NewGetter(m.sagemakerClient)
}

func (m *Manager) Delete(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to userprofile.Options")
	}

	domainID := options.DomainID
	userProfileName := options.Name

	if domainID == "" {
		userProfile, err := m.userProfileGetter.GetUserProfileByName(userProfileName)

		if err != nil {
			var rnfe *resource.NotFoundByNameError
			if errors.As(err, &rnfe) {
				fmt.Printf("SageMaker User Profile with name %q does not exist\n", userProfileName)
				return nil
			}
			return err
		}
		domainID = awssdk.ToString(userProfile.DomainId)
	}

	err := m.sagemakerClient.DeleteUserProfile(domainID, userProfileName)
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}
	fmt.Printf("SageMaker User Profile %q with Domain Id %q deleted\n", userProfileName, domainID)

	return nil
}
