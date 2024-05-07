package ami

import (
	"fmt"
	"os"
	"sort"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	ec2Client *aws.EC2Client
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{ec2Client}
}

func (g *Getter) Init() {
	if g.ec2Client == nil {
		g.ec2Client = aws.NewEC2Client()
	}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	o, ok := options.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to ami.Options")
	}

	filters := []types.Filter{}
	imageIds := []string{}

	if o.NameFilter != "" {
		filters = append(filters, aws.NewEC2NameFilter(o.NameFilter))
	}

	if id != "" {
		filters = []types.Filter{}
		o.Owners = []string{}
		imageIds = append(imageIds, id)
	}

	amis, err := g.ec2Client.DescribeImages(filters, imageIds, o.Owners)
	if err != nil {
		return err
	}

	// Show newer AMIs at the end of the list
	sort.Slice(amis, func(i, j int) bool {
		return awssdk.ToString(amis[i].CreationDate) < awssdk.ToString(amis[j].CreationDate)
	})

	return output.Print(os.Stdout, NewPrinter(amis))
}
