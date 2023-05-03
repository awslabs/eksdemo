package acm_certificate

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	acmClient *aws.ACMClient
}

func NewGetter(acmClient *aws.ACMClient) *Getter {
	return &Getter{acmClient}
}

func (g *Getter) Init() {
	if g.acmClient == nil {
		g.acmClient = aws.NewACMClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var certs []*types.CertificateDetail
	var err error

	if name != "" {
		certs, err = g.GetAllCertsStartingWithName(name)
	} else {
		certs, err = g.GetAllCerts()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(certs))
}

func (g *Getter) GetAllCerts() ([]*types.CertificateDetail, error) {
	certSummaries, err := g.acmClient.ListCertificates()
	if err != nil {
		return nil, err
	}
	certs := make([]*types.CertificateDetail, 0, len(certSummaries))

	for _, summary := range certSummaries {
		cert, err := g.acmClient.DescribeCertificate(awssdk.ToString(summary.CertificateArn))
		if err != nil {
			return nil, err
		}
		certs = append(certs, cert)
	}

	return certs, nil
}

func (g *Getter) GetCert(arn string) (*types.CertificateDetail, error) {
	return g.acmClient.DescribeCertificate(arn)
}

func (g *Getter) GetOneCertStartingWithName(name string) (*types.CertificateDetail, error) {
	certs, err := g.GetAllCertsStartingWithName(name)
	if err != nil {
		return nil, err
	}

	if len(certs) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("certificate name %q not found", name))
	}

	if len(certs) > 1 {
		return nil, fmt.Errorf("multiple certificates found starting with: %s", name)
	}

	return certs[0], nil
}

func (g *Getter) GetAllCertsStartingWithName(name string) ([]*types.CertificateDetail, error) {
	certSummaries, err := g.acmClient.ListCertificates()
	if err != nil {
		return nil, err
	}

	certs := []*types.CertificateDetail{}
	n := strings.ToLower(name)

	for _, summary := range certSummaries {
		if strings.HasPrefix(strings.ToLower(awssdk.ToString(summary.DomainName)), n) {
			cert, err := g.acmClient.DescribeCertificate(awssdk.ToString(summary.CertificateArn))
			if err != nil {
				return nil, err
			}
			certs = append(certs, cert)
		}
	}

	return certs, nil
}
