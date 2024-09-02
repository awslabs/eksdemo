package domain

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

func (g *Getter) Get(domainName string, output printer.Output, o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to domain.Options")
	}

	var domain *sagemaker.DescribeDomainOutput
	var domains []*sagemaker.DescribeDomainOutput
	var err error

	switch {
	case domainName != "":
		domain, err = g.GetDomainByName(domainName)
		domains = []*sagemaker.DescribeDomainOutput{domain}
	case options.DomainID != "":
		domain, err = g.GetDomainByID(options.DomainID)
		domains = []*sagemaker.DescribeDomainOutput{domain}
	default:
		domains, err = g.GetAllDomains()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(domains))
}

func (g *Getter) GetAllDomains() ([]*sagemaker.DescribeDomainOutput, error) {
	domainDetails, err := g.sagemakerClient.ListDomains()
	if err != nil {
		return nil, err
	}

	domains := make([]*sagemaker.DescribeDomainOutput, 0, len(domainDetails))

	for _, dd := range domainDetails {
		result, err := g.sagemakerClient.DescribeDomain(awssdk.ToString(dd.DomainId))
		if err != nil {
			return nil, err
		}
		domains = append(domains, result)
	}

	return domains, nil
}

func (g *Getter) GetDomainByID(domainID string) (*sagemaker.DescribeDomainOutput, error) {
	domain, err := g.sagemakerClient.DescribeDomain(domainID)
	if err != nil {
		return nil, aws.FormatError(err)
	}

	return domain, nil
}

func (g *Getter) GetDomainByName(domainName string) (*sagemaker.DescribeDomainOutput, error) {
	domainDetails, err := g.sagemakerClient.ListDomains()
	if err != nil {
		return nil, err
	}

	found := []types.DomainDetails{}

	for _, dd := range domainDetails {
		if strings.EqualFold(domainName, awssdk.ToString(dd.DomainName)) {
			found = append(found, dd)
		}
	}

	if len(found) == 0 {
		return nil, &resource.NotFoundByNameError{Type: "sagemaker-domain", Name: domainName}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple sagemaker domains found with name: %s", domainName)
	}

	domain, err := g.sagemakerClient.DescribeDomain(awssdk.ToString(found[0].DomainId))
	if err != nil {
		return nil, err
	}

	return domain, nil
}
