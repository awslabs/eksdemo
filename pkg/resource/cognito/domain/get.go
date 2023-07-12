package domain

import (
	"os"

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

func (g *Getter) Get(domain string, output printer.Output, _ resource.Options) error {
	userPoolDomain, err := g.cognitoClient.DescribeUserPooDomainl(domain)

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(userPoolDomain))
}
