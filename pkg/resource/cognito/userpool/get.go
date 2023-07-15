package userpool

import (
	"errors"
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	cognitoClient *aws.CognitoUserPoolClient
}

func NewGetter(cognitoClient *aws.CognitoUserPoolClient) *Getter {
	return &Getter{cognitoClient}
}

func (g *Getter) Init() {
	if g.cognitoClient == nil {
		g.cognitoClient = aws.NewCognitoUserPoolClient()
	}
}

func (g *Getter) Get(id string, output printer.Output, _ resource.Options) error {
	var userpool *types.UserPoolType
	var userpools []*types.UserPoolType
	var err error

	if id != "" {
		userpool, err = g.GetUserPoolByID(id)
		userpools = []*types.UserPoolType{userpool}
	} else {
		userpools, err = g.GetAllUserPools()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(userpools))
}

func (g *Getter) GetAllUserPools() ([]*types.UserPoolType, error) {
	upDescriptions, err := g.cognitoClient.ListUserPools()
	if err != nil {
		return nil, err
	}

	userpools := make([]*types.UserPoolType, 0, len(upDescriptions))

	for _, upd := range upDescriptions {
		result, err := g.cognitoClient.DescribeUserPool(awssdk.ToString(upd.Id))
		if err != nil {
			return nil, err
		}
		userpools = append(userpools, result)
	}

	return userpools, nil
}

func (g *Getter) GetUserPoolByID(id string) (*types.UserPoolType, error) {
	userpool, err := g.cognitoClient.DescribeUserPool(id)

	var rnfe *types.ResourceNotFoundException
	if err != nil && errors.As(err, &rnfe) {
		return nil, &resource.NotFoundByIDError{Type: "user-pool", ID: id}
	}

	return userpool, err
}

func (g *Getter) GetUserPoolByName(name string) (*types.UserPoolType, error) {
	userpools, err := g.GetAllUserPools()
	if err != nil {
		return nil, err
	}

	found := []*types.UserPoolType{}

	for _, up := range userpools {
		if strings.EqualFold(name, awssdk.ToString(up.Name)) {
			found = append(found, up)
		}
	}

	if len(found) == 0 {
		return nil, &resource.NotFoundByNameError{Type: "cognito-user-pool", Name: name}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple cognito user pools found with name: %s", name)
	}

	return found[0], nil
}
