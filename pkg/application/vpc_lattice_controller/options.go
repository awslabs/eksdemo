package vpc_lattice_controller

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type GatwayApiControllerOptions struct {
	application.ApplicationOptions

	Replicas               int
	VpcLatticePrefixListId string
}

func newOptions() (options *GatwayApiControllerOptions, flags cmd.Flags) {
	options = &GatwayApiControllerOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "vpc-lattice",
			ServiceAccount: "gateway-api-controller",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v0.0.16",
				Latest:        "v0.0.16",
				PreviousChart: "v0.0.12",
				Previous:      "v0.0.12",
			},
		},
		Replicas: 1,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
	}

	return
}

func (o *GatwayApiControllerOptions) PreDependencies(application.Action) error {
	pl, err := aws.NewEC2Client().DescribeManagedPrefixLists([]types.Filter{
		{
			Name:   awssdk.String("prefix-list-name"),
			Values: []string{fmt.Sprintf("com.amazonaws.%s.vpc-lattice", o.Region)},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to lookup VPC Lattice Managed Prefix List: %w", err)
	}

	if len(pl) == 0 {
		return fmt.Errorf("failed to lookup VPC Lattice Managed Prefix List: no Prefix List found")
	}

	o.VpcLatticePrefixListId = awssdk.ToString(pl[0].PrefixListId)

	return err
}
