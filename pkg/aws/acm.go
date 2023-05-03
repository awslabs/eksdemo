package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/aws/smithy-go/middleware"
	smithytime "github.com/aws/smithy-go/time"
	smithywaiter "github.com/aws/smithy-go/waiter"
	"github.com/jmespath/go-jmespath"
)

type ACMClient struct {
	*acm.Client
}

func NewACMClient() *ACMClient {
	return &ACMClient{acm.NewFromConfig(GetConfig())}
}

func (c *ACMClient) DeleteCertificate(arn string) error {
	_, err := c.Client.DeleteCertificate(context.Background(), &acm.DeleteCertificateInput{
		CertificateArn: aws.String(arn),
	})

	return err
}

func (c *ACMClient) DescribeCertificate(arn string) (*types.CertificateDetail, error) {
	cert, err := c.Client.DescribeCertificate(context.Background(), &acm.DescribeCertificateInput{
		CertificateArn: aws.String(arn),
	})

	if err != nil {
		return nil, err
	}

	return cert.Certificate, nil
}

func (c *ACMClient) ListCertificates() ([]types.CertificateSummary, error) {
	certs := []types.CertificateSummary{}
	pageNum := 0

	paginator := acm.NewListCertificatesPaginator(c.Client, &acm.ListCertificatesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		certs = append(certs, out.CertificateSummaryList...)
		pageNum++
	}

	return certs, nil
}

func (c *ACMClient) RequestCertificate(fqdn string, sans []string) (string, error) {
	input := &acm.RequestCertificateInput{
		DomainName:       aws.String(fqdn),
		ValidationMethod: types.ValidationMethodDns,
	}

	if len(sans) > 0 {
		input.SubjectAlternativeNames = sans
	}

	out, err := c.Client.RequestCertificate(context.Background(), input)
	if err != nil {
		return "", err
	}

	err = newCertificateValidationMetadataWaiter(c.Client).Wait(context.Background(),
		&acm.DescribeCertificateInput{CertificateArn: out.CertificateArn},
		1*time.Minute,
	)

	return aws.ToString(out.CertificateArn), err
}

type CertificateValidationMetadataWaiterOptions acm.CertificateValidatedWaiterOptions

// CertificateValidationMetadataWaiter defines the waiters for CertificateValidationMetadata
type CertificateValidationMetadataWaiter struct {
	client acm.DescribeCertificateAPIClient

	options CertificateValidationMetadataWaiterOptions
}

// newCertificateValidationMetadataWaiter constructs a CertificateValidationMetadataWaiter.
func newCertificateValidationMetadataWaiter(client acm.DescribeCertificateAPIClient, optFns ...func(*CertificateValidationMetadataWaiterOptions)) *CertificateValidationMetadataWaiter {
	options := CertificateValidationMetadataWaiterOptions{}
	options.APIOptions = append(options.APIOptions, WaiterLogger{}.AddLogger)
	options.MinDelay = 2 * time.Second
	options.MaxDelay = 120 * time.Second
	options.Retryable = certificateHasValidationMetadataRetryable

	for _, fn := range optFns {
		fn(&options)
	}
	return &CertificateValidationMetadataWaiter{
		client:  client,
		options: options,
	}
}

// Wait calls the waiter function for CertificateValidationMetadata waiter. The maxWaitDur
// is the maximum wait duration the waiter will wait. The maxWaitDur is required
// and must be greater than zero.
func (w *CertificateValidationMetadataWaiter) Wait(ctx context.Context, params *acm.DescribeCertificateInput, maxWaitDur time.Duration, optFns ...func(*CertificateValidationMetadataWaiterOptions)) error {
	_, err := w.WaitForOutput(ctx, params, maxWaitDur, optFns...)
	return err
}

// WaitForOutput calls the waiter function for CertificateValidationMetadata waiter and
// returns the output of the successful operation. The maxWaitDur is the maximum
// wait duration the waiter will wait. The maxWaitDur is required and must be
// greater than zero.
func (w *CertificateValidationMetadataWaiter) WaitForOutput(ctx context.Context, params *acm.DescribeCertificateInput, maxWaitDur time.Duration, optFns ...func(*CertificateValidationMetadataWaiterOptions)) (*acm.DescribeCertificateOutput, error) {
	if maxWaitDur <= 0 {
		return nil, fmt.Errorf("maximum wait time for waiter must be greater than zero")
	}

	options := w.options
	for _, fn := range optFns {
		fn(&options)
	}

	if options.MaxDelay <= 0 {
		options.MaxDelay = 120 * time.Second
	}

	if options.MinDelay > options.MaxDelay {
		return nil, fmt.Errorf("minimum waiter delay %v must be lesser than or equal to maximum waiter delay of %v", options.MinDelay, options.MaxDelay)
	}

	ctx, cancelFn := context.WithTimeout(ctx, maxWaitDur)
	defer cancelFn()

	logger := smithywaiter.Logger{}
	remainingTime := maxWaitDur

	var attempt int64
	for {

		attempt++
		apiOptions := options.APIOptions
		start := time.Now()

		if options.LogWaitAttempts {
			logger.Attempt = attempt
			apiOptions = append([]func(*middleware.Stack) error{}, options.APIOptions...)
			apiOptions = append(apiOptions, logger.AddLogger)
		}

		out, err := w.client.DescribeCertificate(ctx, params, func(o *acm.Options) {
			o.APIOptions = append(o.APIOptions, apiOptions...)
		})

		retryable, err := options.Retryable(ctx, params, out, err)
		if err != nil {
			return nil, err
		}
		if !retryable {
			return out, nil
		}

		remainingTime -= time.Since(start)
		if remainingTime < options.MinDelay || remainingTime <= 0 {
			break
		}

		// compute exponential backoff between waiter retries
		delay, err := smithywaiter.ComputeDelay(
			attempt, options.MinDelay, options.MaxDelay, remainingTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error computing waiter delay, %w", err)
		}

		remainingTime -= delay
		// sleep for the delay amount before invoking a request
		if err := smithytime.SleepWithContext(ctx, delay); err != nil {
			return nil, fmt.Errorf("request cancelled while waiting, %w", err)
		}
	}
	return nil, fmt.Errorf("exceeded max wait time for CertificateValidationMetadata waiter")
}

func certificateHasValidationMetadataRetryable(ctx context.Context, input *acm.DescribeCertificateInput, output *acm.DescribeCertificateOutput, err error) (bool, error) {

	if err == nil {
		pathValue, err := jmespath.Search("Certificate.DomainValidationOptions[].ResourceRecord.Type", output)
		if err != nil {
			return false, fmt.Errorf("error evaluating waiter state: %w", err)
		}

		expectedValue := "CNAME"
		var match = true
		listOfValues, ok := pathValue.([]interface{})
		if !ok {
			return false, fmt.Errorf("waiter comparator expected list got %T", pathValue)
		}

		if len(listOfValues) == 0 {
			match = false
		}
		for _, v := range listOfValues {
			value, ok := v.(types.RecordType)
			if !ok {
				return false, fmt.Errorf("waiter comparator expected types.RecordType value, got %T", pathValue)
			}

			if string(value) != expectedValue {
				match = false
			}
		}

		if match {
			return false, nil
		}
	}

	if err == nil {
		pathValue, err := jmespath.Search("Certificate.Status", output)
		if err != nil {
			return false, fmt.Errorf("error evaluating waiter state: %w", err)
		}

		expectedValue := "FAILED"
		value, ok := pathValue.(types.CertificateStatus)
		if !ok {
			return false, fmt.Errorf("waiter comparator expected types.CertificateStatus value, got %T", pathValue)
		}

		if string(value) == expectedValue {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	if err != nil {
		var errorType *types.ResourceNotFoundException
		if errors.As(err, &errorType) {
			return false, fmt.Errorf("waiter state transitioned to Failure")
		}
	}

	return true, nil
}
