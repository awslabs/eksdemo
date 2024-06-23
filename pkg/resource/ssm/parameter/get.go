package parameter

import (
	"errors"
	"fmt"
	"os"
	"sort"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	ssmClient *aws.SSMClient
}

func NewGetter(ssmClient *aws.SSMClient) *Getter {
	return &Getter{ssmClient}
}

func (g *Getter) Init() {
	if g.ssmClient == nil {
		g.ssmClient = aws.NewSSMClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to parameter.Options")
	}

	switch {
	case options.Path != "":
		params, err := g.GetByPath(options.Path)
		if err != nil {
			return err
		}
		// Show recently updated Parameters at the end of the list
		sort.Slice(params, func(i, j int) bool {
			return params[i].LastModifiedDate.Before(awssdk.ToTime(params[j].LastModifiedDate))
		})
		return output.Print(os.Stdout, NewPrinter(params))

	case name != "":
		param, err := g.GetByName(name)
		var rnfe *types.ParameterNotFound
		if err != nil {
			if errors.As(err, &rnfe) {
				return fmt.Errorf("ssm-parameter name %q not found", name)
			}
			return err
		}
		return output.Print(os.Stdout, NewPrinter([]types.Parameter{*param}))

	default:
		params, err := g.ssmClient.DescribeParameters()
		if err != nil {
			return err
		}
		return output.Print(os.Stdout, NewMetadataPrinter(params))
	}

}

func (g *Getter) GetByName(name string) (*types.Parameter, error) {
	param, err := g.ssmClient.GetParameter(name)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (g *Getter) GetByPath(pathOrName string) ([]types.Parameter, error) {
	params, err := g.ssmClient.GetParametersByPath(pathOrName)
	if err != nil {
		return nil, err
	}
	return params, nil
}
