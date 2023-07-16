package client

import (
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

func (g *Getter) Get(name string, output printer.Output, o resource.Options) error {
	options, ok := o.(*Options)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to client.Options")
	}

	var appClient *types.UserPoolClientType
	var appClients []*types.UserPoolClientType
	var err error

	switch {
	case name != "":
		appClient, err = g.GetAppClientByName(name, options.UserPoolID)
		appClients = []*types.UserPoolClientType{appClient}
	case options.AppClientID != "":
		appClient, err = g.GetAppClientByID(options.AppClientID, options.UserPoolID)
		appClients = []*types.UserPoolClientType{appClient}
	default:
		appClients, err = g.GetAllAppClients(options.UserPoolID)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(appClients))
}

func (g *Getter) GetAllAppClients(userPoolID string) ([]*types.UserPoolClientType, error) {
	appClientDescriptions, err := g.cognitoClient.ListUserPoolClients(userPoolID)
	if err != nil {
		return nil, err
	}

	appClients := make([]*types.UserPoolClientType, 0, len(appClientDescriptions))

	for _, acd := range appClientDescriptions {
		result, err := g.cognitoClient.DescribeUserPoolClient(awssdk.ToString(acd.ClientId), userPoolID)
		if err != nil {
			return nil, err
		}
		appClients = append(appClients, result)
	}

	return appClients, nil
}

func (g *Getter) GetAppClientByID(appClientID, userPoolID string) (*types.UserPoolClientType, error) {
	appClient, err := g.cognitoClient.DescribeUserPoolClient(appClientID, userPoolID)
	if err != nil {
		return nil, err
	}

	return appClient, nil
}

func (g *Getter) GetAppClientByName(name, userPoolID string) (*types.UserPoolClientType, error) {
	appClients, err := g.GetAllAppClients(userPoolID)
	if err != nil {
		return nil, err
	}

	found := []*types.UserPoolClientType{}

	for _, ac := range appClients {
		if strings.EqualFold(name, awssdk.ToString(ac.ClientName)) {
			found = append(found, ac)
		}
	}

	if len(found) == 0 {
		return nil, &resource.NotFoundByNameError{Type: "cognito-app-client", Name: name}
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple cognito app clients found with name: %s", name)
	}

	return found[0], nil
}
