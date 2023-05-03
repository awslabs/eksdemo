package application

import (
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/helm"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	resource.EmptyInit
}

func NewGetter(ec2Client *aws.EC2Client) *Getter {
	return &Getter{}
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	releases, err := helm.List(options.Common().KubeContext)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(releases))
}
