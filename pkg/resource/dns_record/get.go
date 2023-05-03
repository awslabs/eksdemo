package dns_record

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/hosted_zone"
)

type Getter struct {
	route53Client *aws.Route53Client
	zoneGetter    *hosted_zone.Getter
}

func NewGetter(route53Client *aws.Route53Client) *Getter {
	return &Getter{route53Client, hosted_zone.NewGetter(route53Client)}
}

func (g *Getter) Init() {
	if g.route53Client == nil {
		g.route53Client = aws.NewRoute53Client()
	}
	g.zoneGetter = hosted_zone.NewGetter(g.route53Client)
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := g.zoneGetter.GetZoneByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	filterTypes := map[string]bool{}
	for _, f := range dnsOptions.Filter {
		filterTypes[f] = true
	}

	recordSets, err := g.GetRecordsWithFilter(name, awssdk.ToString(zone.Id), filterTypes)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(recordSets))
}

func (g *Getter) GetRecords(name, zoneId string) ([]types.ResourceRecordSet, error) {
	return g.GetRecordsWithFilter(name, zoneId, map[string]bool{})
}

func (g *Getter) GetRecordsWithFilter(name, zoneId string, filterTypes map[string]bool) ([]types.ResourceRecordSet, error) {
	recordSets, err := g.route53Client.ListResourceRecordSets(zoneId)
	if err != nil {
		return nil, err
	}

	if name != "" {
		n := strings.ToLower(name) + "."
		filtered := []types.ResourceRecordSet{}
		for _, rs := range recordSets {
			if n == strings.ToLower(awssdk.ToString(rs.Name)) {
				filtered = append(filtered, rs)
			}
		}
		recordSets = filtered
	}

	if len(filterTypes) > 0 {
		filtered := make([]types.ResourceRecordSet, 0, len(recordSets))
		for _, rs := range recordSets {
			if filterTypes[string(rs.Type)] {
				filtered = append(filtered, rs)
			}
		}
		recordSets = filtered
	}

	return recordSets, nil
}
