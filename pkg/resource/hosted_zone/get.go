package hosted_zone

import (
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	route53Client *aws.Route53Client
}

func NewGetter(route53Client *aws.Route53Client) *Getter {
	return &Getter{route53Client}
}

func (g *Getter) Init() {
	if g.route53Client == nil {
		g.route53Client = aws.NewRoute53Client()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var err error
	var zone types.HostedZone
	var zones []types.HostedZone

	if name != "" {
		zone, err = g.GetZoneByName(name)
		zones = []types.HostedZone{zone}
	} else {
		zones, err = g.GetAllZones()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(zones))
}

func (g *Getter) GetAllZones() ([]types.HostedZone, error) {
	return g.route53Client.ListHostedZones()
}

func (g *Getter) GetZoneByName(name string) (types.HostedZone, error) {
	zone, err := g.route53Client.ListHostedZonesByName(name)
	if err != nil {
		return types.HostedZone{}, err
	}

	if len(zone) == 0 {
		return types.HostedZone{}, fmt.Errorf("hosted-zone %q not found", name)
	}

	z := zone[0]

	if strings.ToLower(awssdk.ToString(z.Name)) != name+"." {
		return types.HostedZone{}, fmt.Errorf("hosted-zone %q not found", name)
	}

	return z, nil
}
