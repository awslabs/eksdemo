package session

import (
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	resource.EmptyInit
}

func (g *Getter) Get(id string, output printer.Output, options resource.Options) error {
	sessOptions, ok := options.(*SessionOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SessionOptions")
	}

	state := "Active"
	if sessOptions.History {
		state = "History"
	}

	sessions, err := aws.NewSSMClient().DescribeSessions(id, state)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(sessions))
}
