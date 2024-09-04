package filesystem

import (
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	efsClient *aws.EFSClient
}

func NewGetter(efsClient *aws.EFSClient) *Getter {
	return &Getter{efsClient}
}

func (g *Getter) Init() {
	if g.efsClient == nil {
		g.efsClient = aws.NewEFSClient()
	}
}

func (g *Getter) Get(id string, output printer.Output, _ resource.Options) error {
	fileSystems, err := g.efsClient.DescribeFileSystems(id)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(fileSystems))
}
