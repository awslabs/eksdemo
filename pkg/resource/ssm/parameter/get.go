package parameter

import (
	"errors"
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

func (g *Getter) Get(pathOrName string, output printer.Output, _ resource.Options) error {
	params, err := g.GetByPathOrName(pathOrName)
	if err != nil {
		return err
	}

	// Show recently updated Parameters at the end of the list
	sort.Slice(params, func(i, j int) bool {
		return params[i].LastModifiedDate.Before(awssdk.ToTime(params[j].LastModifiedDate))
	})

	return output.Print(os.Stdout, NewPrinter(params))
}

func (g *Getter) GetByPathOrName(pathOrName string) ([]types.Parameter, error) {
	params, err := g.ssmClient.GetParametersByPath(pathOrName)
	if err != nil {
		return nil, err
	}

	if len(params) > 0 {
		return params, nil
	}

	param, err := g.ssmClient.GetParameter(pathOrName)

	// Return all errors except NotFound
	var rnfe *types.ParameterNotFound
	if err != nil && !errors.As(err, &rnfe) {
		return nil, err
	}

	if param != nil {
		return []types.Parameter{*param}, nil
	}

	return nil, nil
}
