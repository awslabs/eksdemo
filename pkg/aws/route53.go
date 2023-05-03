package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53Client struct {
	*route53.Client
}

func NewRoute53Client() *Route53Client {
	return &Route53Client{route53.NewFromConfig(GetConfig())}
}

func (c *Route53Client) ChangeResourceRecordSets(changeBatch *types.ChangeBatch, zoneId string) error {
	_, err := c.Client.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: aws.String(zoneId),
	})

	return err
}

func (c *Route53Client) GetHostedZone(zoneId string) (*route53.GetHostedZoneOutput, error) {
	zone, err := c.Client.GetHostedZone(context.Background(), &route53.GetHostedZoneInput{
		Id: aws.String(zoneId),
	})

	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (c *Route53Client) ListHostedZones() ([]types.HostedZone, error) {
	zones := []types.HostedZone{}
	pageNum := 0

	paginator := route53.NewListHostedZonesPaginator(c.Client, &route53.ListHostedZonesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		zones = append(zones, out.HostedZones...)
		pageNum++
	}

	return zones, nil
}

func (c *Route53Client) ListHostedZonesByName(name string) ([]types.HostedZone, error) {
	zones, err := c.Client.ListHostedZonesByName(context.Background(), &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(name),
	})

	if err != nil {
		return nil, err
	}

	return zones.HostedZones, nil
}

func (c *Route53Client) ListResourceRecordSets(zoneId string) ([]types.ResourceRecordSet, error) {
	recordSets := []types.ResourceRecordSet{}
	pageNum := 0

	paginator := NewListResourceRecordSetsPaginator(c.Client, &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		recordSets = append(recordSets, out.ResourceRecordSets...)
		pageNum++
	}

	return recordSets, nil
}

// ListResourceRecordSetsAPIClient is a client that implements the ListResourceRecordSets
// operation.
type ListResourceRecordSetsAPIClient interface {
	ListResourceRecordSets(context.Context, *route53.ListResourceRecordSetsInput, ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error)
}

var _ ListResourceRecordSetsAPIClient = (*route53.Client)(nil)

// ListResourceRecordSetsPaginatorOptions is the paginator options for ListResourceRecordSets
type ListResourceRecordSetsPaginatorOptions struct {
	// (Optional) The maximum number of resource record sets that you want Amazon Route 53 to
	// return. If you have more than maxitems resource record sets, the value of IsTruncated in
	// the response is true, and the value of NextMarker is the resource record set ID of the
	// first esource record set that Route 53 will return if you submit another request.
	Limit int32

	// Set to true if pagination should stop if the service returns a pagination token
	// that matches the most recent token provided to the service.
	StopOnDuplicateToken bool
}

// ListResourceRecordSetsPaginator is a paginator for ListResourceRecordSets
type ListResourceRecordSetsPaginator struct {
	options              ListResourceRecordSetsPaginatorOptions
	client               ListResourceRecordSetsAPIClient
	params               *route53.ListResourceRecordSetsInput
	nextRecordIdentifier *string
	nextRecordName       *string
	nextRecordType       types.RRType
	firstPage            bool
}

// NewListResourceRecordSetsPaginator returns a new ListResourceRecordSetsPaginator
func NewListResourceRecordSetsPaginator(client ListResourceRecordSetsAPIClient, params *route53.ListResourceRecordSetsInput, optFns ...func(*ListResourceRecordSetsPaginatorOptions)) *ListResourceRecordSetsPaginator {
	if params == nil {
		params = &route53.ListResourceRecordSetsInput{}
	}

	options := ListResourceRecordSetsPaginatorOptions{}
	if params.MaxItems != nil {
		options.Limit = *params.MaxItems
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &ListResourceRecordSetsPaginator{
		options:              options,
		client:               client,
		params:               params,
		firstPage:            true,
		nextRecordIdentifier: params.StartRecordIdentifier,
		nextRecordName:       params.StartRecordName,
		nextRecordType:       params.StartRecordType,
	}
}

// HasMorePages returns a boolean indicating whether more pages are available
func (p *ListResourceRecordSetsPaginator) HasMorePages() bool {
	return p.firstPage || (p.nextRecordIdentifier != nil && len(*p.nextRecordIdentifier) != 0)
}

// NextPage retrieves the next ListResourceRecordSets page.
func (p *ListResourceRecordSetsPaginator) NextPage(ctx context.Context, optFns ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error) {
	if !p.HasMorePages() {
		return nil, fmt.Errorf("no more pages available")
	}

	params := *p.params
	params.StartRecordIdentifier = p.nextRecordIdentifier
	params.StartRecordName = p.nextRecordName
	params.StartRecordType = p.nextRecordType

	var limit *int32
	if p.options.Limit > 0 {
		limit = &p.options.Limit
	}
	params.MaxItems = limit

	result, err := p.client.ListResourceRecordSets(ctx, &params, optFns...)
	if err != nil {
		return nil, err
	}
	p.firstPage = false

	prevRecordIdentifier := p.nextRecordIdentifier
	p.nextRecordIdentifier = result.NextRecordIdentifier
	p.nextRecordName = result.NextRecordName
	p.nextRecordType = result.NextRecordType

	if p.options.StopOnDuplicateToken &&
		prevRecordIdentifier != nil &&
		p.nextRecordIdentifier != nil &&
		*prevRecordIdentifier == *p.nextRecordIdentifier {
		p.nextRecordIdentifier = nil
	}

	return result, nil
}
