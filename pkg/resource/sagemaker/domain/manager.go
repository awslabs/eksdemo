package domain

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
	sagemakerClient *aws.SageMakerClient
	domainGetter    *Getter
}

func (m *Manager) Init() {
	if m.sagemakerClient == nil {
		m.sagemakerClient = aws.NewSageMakerClient()
	}
	m.domainGetter = NewGetter(m.sagemakerClient)
}

func (m *Manager) Delete(o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to domain.Options")
	}

	domainID := options.DomainID
	domainName := options.Name

	if domainID == "" {
		domain, err := m.domainGetter.GetDomainByName(domainName)

		if err != nil {
			var rnfe *resource.NotFoundByNameError
			if errors.As(err, &rnfe) {
				fmt.Printf("SageMaker Domain with name %q does not exist\n", domainName)
				return nil
			}
			return err
		}

		domainID = awssdk.ToString(domain.DomainId)
	}

	err := m.sagemakerClient.DeleteDomain(domainID)
	if err != nil {
		return aws.FormatErrorAsMessageOnly(err)
	}
	fmt.Printf("SageMaker Domain Id %q deleted\n", domainID)

	return nil
}
