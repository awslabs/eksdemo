package node

import (
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(instanceId string, output printer.Output, options resource.Options) error {
	nodes, err := aws.NewSSMClient().DescribeInstanceInformation(instanceId)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(nodes))
}
