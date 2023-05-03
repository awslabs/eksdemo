package prefix_list

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type PrefixList struct {
	PrefixList types.ManagedPrefixList
	Entries    []types.PrefixListEntry
}

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
	var filters []types.Filter
	var prefixLists []PrefixList
	var err error

	if id != "" {
		prefixLists, err = g.GetByPrefixListId(id)
	} else {
		prefixLists, err = g.GetAll(filters)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(prefixLists))
}

func (g *Getter) GetAll(filters []types.Filter) ([]PrefixList, error) {
	pls, err := g.ec2Client.DescribeManagedPrefixLists(filters)
	if err != nil {
		return nil, err
	}

	prefixLists := make([]PrefixList, 0, len(pls))
	for _, pl := range pls {
		entries, err := g.ec2Client.GetManagedPrefixListEntries(awssdk.ToString(pl.PrefixListId))
		if err != nil {
			return nil, err
		}
		prefixLists = append(prefixLists, PrefixList{pl, entries})

	}
	return prefixLists, nil
}

func (g *Getter) GetByPrefixListId(prefixListId string) ([]PrefixList, error) {
	pl, err := g.GetAll([]types.Filter{aws.NewEC2PrefixListFilter(prefixListId)})
	if err != nil {
		return nil, err
	}

	if len(pl) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("prefix-list %q not found", prefixListId))
	}

	return pl, nil
}
