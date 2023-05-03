package cloudtrail_trail

import (
	"errors"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Trail struct {
	Trail       *types.Trail
	TrailStatus *cloudtrail.GetTrailStatusOutput
}

type Getter struct {
	cloudtrailClient *aws.CloudtrailClient
}

func NewGetter(cloudtrailClient *aws.CloudtrailClient) *Getter {
	return &Getter{cloudtrailClient}
}

func (g *Getter) Init() {
	if g.cloudtrailClient == nil {
		g.cloudtrailClient = aws.NewCloudtrailClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var trail Trail
	var trails []Trail
	var err error

	if name != "" {
		trail, err = g.GetTrailByName(name)
		trails = []Trail{trail}
	} else {
		trails, err = g.GetAllTrails()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(trails))
}

func (g *Getter) GetAllTrails() ([]Trail, error) {
	trailInfos, err := g.cloudtrailClient.ListTrails()
	if err != nil {
		return nil, err
	}

	trails := make([]Trail, 0, len(trailInfos))

	for _, t := range trailInfos {
		trail, err := g.cloudtrailClient.GetTrail(awssdk.ToString(t.TrailARN))
		if err != nil {
			return nil, err
		}

		status, err := g.cloudtrailClient.GetTrailStatus(awssdk.ToString(t.TrailARN))
		if err != nil {
			return nil, err
		}
		trails = append(trails, Trail{trail, status})
	}
	return trails, nil
}

func (g *Getter) GetTrailByName(name string) (Trail, error) {
	trail, err := g.cloudtrailClient.GetTrail(name)

	var tnfe *types.TrailNotFoundException
	if err != nil {
		if errors.As(err, &tnfe) {
			return Trail{}, resource.NotFoundError(fmt.Sprintf("cloudtrail-trail %q not found in region %q",
				name, g.cloudtrailClient.Region))
		}
		return Trail{}, err
	}

	status, err := g.cloudtrailClient.GetTrailStatus(name)

	return Trail{trail, status}, err
}
