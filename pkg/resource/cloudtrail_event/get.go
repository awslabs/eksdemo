package cloudtrail_event

import (
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

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

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	eventOptions, ok := options.(*CloudtrailEventOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to CloudtrailEventOptions")
	}

	filters := []types.LookupAttribute{}

	switch {
	case id != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyEventId,
			AttributeValue: awssdk.String(id),
		})

	case eventOptions.Name != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyEventName,
			AttributeValue: awssdk.String(eventOptions.Name),
		})

	case eventOptions.ResourceName != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyResourceName,
			AttributeValue: awssdk.String(eventOptions.ResourceName),
		})

	case eventOptions.ResourceType != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyResourceType,
			AttributeValue: awssdk.String(eventOptions.ResourceType),
		})

	case eventOptions.Source != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyEventSource,
			AttributeValue: awssdk.String(eventOptions.Source),
		})

	case eventOptions.Username != "":
		filters = append(filters, types.LookupAttribute{
			AttributeKey:   types.LookupAttributeKeyUsername,
			AttributeValue: awssdk.String(eventOptions.Username),
		})
	}

	events, err := g.cloudtrailClient.LookupEvents(eventOptions.Insights, filters)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(events))
}
