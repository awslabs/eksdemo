package dns_record

import (
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/hosted_zone"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun          bool
	dnsRecordGetter *Getter
	route53Client   *aws.Route53Client
	zoneGetter      *hosted_zone.Getter
}

func (m *Manager) Init() {
	if m.route53Client == nil {
		m.route53Client = aws.NewRoute53Client()
	}
	m.dnsRecordGetter = NewGetter(m.route53Client)
	m.zoneGetter = hosted_zone.NewGetter(m.route53Client)
}

func (m *Manager) Create(options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := m.zoneGetter.GetZoneByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	changeBatch := &types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: types.ChangeActionCreate,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name: awssdk.String(dnsOptions.Name),
					ResourceRecords: []types.ResourceRecord{
						{
							Value: awssdk.String(dnsOptions.Value),
						},
					},
					TTL:  awssdk.Int64(300),
					Type: types.RRType(dnsOptions.Type),
				},
			},
		},
		Comment: awssdk.String("eksdemo create dns-record"),
	}

	if err := m.route53Client.ChangeResourceRecordSets(changeBatch, awssdk.ToString(zone.Id)); err != nil {
		return err
	}
	fmt.Println("Record created successfully")

	return nil
}

func (m *Manager) Delete(options resource.Options) error {
	dnsOptions, ok := options.(*DnsRecordOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to DnsRecordOptions")
	}

	zone, err := m.zoneGetter.GetZoneByName(dnsOptions.ZoneName)
	if err != nil {
		return err
	}

	records, err := m.dnsRecordGetter.GetRecords(dnsOptions.Name, awssdk.ToString(zone.Id))
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found with name %q", dnsOptions.Name)
	}

	if len(records) > 1 && !dnsOptions.AllTypes && !dnsOptions.AllRecords {
		return fmt.Errorf("multiple records found with name %q, use %q flag to delete all records", dnsOptions.Name, "--all")
	}

	changes := make([]types.Change, 0, len(records))

	for i, rec := range records {
		if rec.Type == types.RRTypeNs || rec.Type == types.RRTypeSoa {
			continue
		}

		change := types.Change{
			Action:            types.ChangeActionDelete,
			ResourceRecordSet: &records[i],
		}
		changes = append(changes, change)
		fmt.Printf("Deleting %s record %q...\n", string(rec.Type), strings.TrimSuffix(awssdk.ToString(rec.Name), "."))
	}

	if len(changes) == 0 {
		fmt.Println("No records to delete.")
		return nil
	}

	changeBatch := &types.ChangeBatch{
		Changes: changes,
		Comment: awssdk.String("eksdemo delete dns-record"),
	}

	if err := m.route53Client.ChangeResourceRecordSets(changeBatch, awssdk.ToString(zone.Id)); err != nil {
		return err
	}
	fmt.Println("Record(s) deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
