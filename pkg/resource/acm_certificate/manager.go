package acm_certificate

import (
	"context"
	"fmt"
	"strings"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/hosted_zone"
	"github.com/spf13/cobra"
)

type Manager struct {
	acmClient     *aws.ACMClient
	acmGetter     *Getter
	route53Client *aws.Route53Client
	zoneGetter    *hosted_zone.Getter
}

func (m *Manager) Init() {
	if m.acmClient == nil {
		m.acmClient = aws.NewACMClient()
	}
	if m.route53Client == nil {
		m.route53Client = aws.NewRoute53Client()
	}
	m.acmGetter = NewGetter(m.acmClient)
	m.zoneGetter = hosted_zone.NewGetter(m.route53Client)
}

func (m *Manager) Create(options resource.Options) error {
	certOptions, ok := options.(*CertificateOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to CertificateOptions")
	}

	name := certOptions.Name
	cert, err := m.acmGetter.GetOneCertStartingWithName(name)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); !ok {
			// Return an error if it's anything other than resource not found
			return err
		}
	}

	if cert != nil && strings.EqualFold(awssdk.ToString(cert.DomainName), name) {
		fmt.Printf("Certificate %q already exists\n", name)
		return nil
	}

	fmt.Printf("Creating ACM Certificate request for: %s...", name)
	arn, err := m.acmClient.RequestCertificate(name, certOptions.sans)
	if err != nil {
		return err
	}
	fmt.Printf("done\nCreated ACM Certificate Id: %s\n", arn[strings.LastIndex(arn, "/")+1:])

	cert, err = m.acmGetter.GetCert(arn)
	if err != nil {
		return fmt.Errorf("failed to describe the certificate: %w", err)
	}

	if certOptions.skipValidation {
		return m.outputValidationSteps(cert)
	}

	return m.validate(cert)
}

func (m *Manager) Delete(options resource.Options) error {
	name := options.Common().Name

	cert, err := m.acmGetter.GetOneCertStartingWithName(name)
	if err != nil {
		return err
	}

	err = m.acmClient.DeleteCertificate(awssdk.ToString(cert.CertificateArn))
	if err != nil {
		return err
	}
	fmt.Printf("ACM Certificate Domain name %q deleted\n", awssdk.ToString(cert.DomainName))

	return nil
}

func (m *Manager) SetDryRun() {}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) outputValidationSteps(cert *types.CertificateDetail) error {
	fmt.Println("To validate:")
	for _, dvo := range cert.DomainValidationOptions {
		fmt.Printf("In zone %q, create %q record %q with value %q\n",
			awssdk.ToString(dvo.DomainName),
			string(dvo.ResourceRecord.Type),
			awssdk.ToString(dvo.ResourceRecord.Name),
			awssdk.ToString(dvo.ResourceRecord.Value),
		)
	}
	return nil
}

func (m *Manager) validate(cert *types.CertificateDetail) error {
	zones, err := m.zoneGetter.GetAllZones()
	if err != nil {
		return fmt.Errorf("failed during validation to list hosted zones: %w", err)
	}

	for _, z := range zones {
		changes := []route53types.Change{}
		zoneName := strings.TrimSuffix(awssdk.ToString(z.Name), ".")

		for _, dv := range cert.DomainValidationOptions {
			if strings.HasSuffix(awssdk.ToString(dv.DomainName), zoneName) {
				fmt.Printf("Validating domain %q using hosted zone %q\n", awssdk.ToString(dv.DomainName), zoneName)
				rr := dv.ResourceRecord
				changes = append(changes, createChange(rr.Name, rr.Value, z.Id, rr.Type))
			}
		}

		if len(changes) == 0 {
			continue
		}

		changeBatch := &route53types.ChangeBatch{
			Changes: changes,
			Comment: awssdk.String("certificate validation"),
		}

		if err := m.route53Client.ChangeResourceRecordSets(changeBatch, awssdk.ToString(z.Id)); err != nil {
			return err
		}
	}

	fmt.Printf("Waiting for certificate to be issued...")
	waiter := acm.NewCertificateValidatedWaiter(m.acmClient.Client, func(o *acm.CertificateValidatedWaiterOptions) {
		o.MinDelay = 2 * time.Second
		o.MaxDelay = 10 * time.Second
		o.APIOptions = append(o.APIOptions, aws.WaiterLogger{}.AddLogger)
	})

	err = waiter.Wait(context.Background(),
		&acm.DescribeCertificateInput{CertificateArn: cert.CertificateArn},
		3*time.Minute,
	)
	if err != nil {
		fmt.Println()
		return err
	}
	fmt.Println("done")

	return nil
}

func createChange(name, value, zoneId *string, recType types.RecordType) route53types.Change {
	return route53types.Change{
		Action: route53types.ChangeActionUpsert,
		ResourceRecordSet: &route53types.ResourceRecordSet{
			Name: name,
			ResourceRecords: []route53types.ResourceRecord{
				{
					Value: value,
				},
			},
			TTL:  awssdk.Int64(300),
			Type: route53types.RRType(recType),
		},
	}
}
