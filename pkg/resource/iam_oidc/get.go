package iam_oidc

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	iamClient *aws.IAMClient
}

func NewGetter(iamClient *aws.IAMClient) *Getter {
	return &Getter{iamClient}
}

func (g *Getter) Init() {
	if g.iamClient == nil {
		g.iamClient = aws.NewIAMClient()
	}
}

func (g *Getter) Get(providerUrl string, output printer.Output, options resource.Options) error {
	var err error
	var oidcProvider *iam.GetOpenIDConnectProviderOutput
	var oidcProviders []*iam.GetOpenIDConnectProviderOutput

	if providerUrl != "" {
		oidcProvider, err = g.GetOidcProviderByUrl(providerUrl)
		oidcProviders = []*iam.GetOpenIDConnectProviderOutput{oidcProvider}
	} else if options.Common().Cluster != nil {
		oidcProvider, err = g.GetOidcProviderByCluster(options.Common().Cluster)
		oidcProviders = []*iam.GetOpenIDConnectProviderOutput{oidcProvider}
	} else {
		oidcProviders, err = g.GetAllOidcProviders()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(oidcProviders))
}

func (g *Getter) GetAllOidcProviders() ([]*iam.GetOpenIDConnectProviderOutput, error) {
	providerList, err := g.iamClient.ListOpenIDConnectProviders()
	if err != nil {
		return nil, err
	}
	oidcProviders := make([]*iam.GetOpenIDConnectProviderOutput, 0, len(providerList))

	for _, p := range providerList {
		provider, err := g.iamClient.GetOpenIDConnectProvider(awssdk.ToString(p.Arn))
		if err != nil {
			return nil, err
		}
		oidcProviders = append(oidcProviders, provider)
	}
	return oidcProviders, nil
}

func (g *Getter) GetOidcProviderByCluster(cluster *ekstypes.Cluster) (*iam.GetOpenIDConnectProviderOutput, error) {
	u, err := url.Parse(awssdk.ToString(cluster.Identity.Oidc.Issuer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL when validating options: %w", err)
	}

	oidc, err := g.GetOidcProviderByUrl(u.Hostname() + u.Path)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); ok {
			return nil, fmt.Errorf("cluster %q has no IAM OIDC identity provider configured", awssdk.ToString(cluster.Name))
		}
		return nil, err
	}

	return oidc, nil
}

func (g *Getter) GetOidcProviderByUrl(url string) (*iam.GetOpenIDConnectProviderOutput, error) {
	arn := fmt.Sprintf("arn:%s:iam::%s:oidc-provider/%s", aws.Partition(), aws.AccountId(), url)

	provider, err := g.iamClient.GetOpenIDConnectProvider(arn)
	if err != nil {
		var nsee *types.NoSuchEntityException
		if errors.As(err, &nsee) {
			return nil, resource.NotFoundError(fmt.Sprintf("oidc-provider %q not found", url))
		}
		return nil, err
	}

	return provider, nil
}
