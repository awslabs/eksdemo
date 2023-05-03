package iam_role

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/iam_oidc"
)

type Getter struct {
	iamClient  *aws.IAMClient
	oidcGetter *iam_oidc.Getter
}

func NewGetter(iamClient *aws.IAMClient) *Getter {
	return &Getter{iamClient, iam_oidc.NewGetter(iamClient)}
}

func (g *Getter) Init() {
	if g.iamClient == nil {
		g.iamClient = aws.NewIAMClient()
	}
	g.oidcGetter = iam_oidc.NewGetter(g.iamClient)
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	roleOptions, ok := options.(*IamRoleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to IamRoleOptions")
	}

	var role *types.Role
	var roles []types.Role
	var err error

	if name != "" {
		role, err = g.GetRoleByName(name)
	} else if roleOptions.NameSearch != "" {
		roles, err = g.GetRolesByNameSearch(roleOptions.NameSearch)
	} else if roleOptions.Cluster != nil {
		roles, err = g.GetIrsaRolesForCluster(roleOptions.Cluster, roleOptions.LastUsed)
	} else {
		roles, err = g.GetAllRoles(roleOptions.LastUsed)
	}

	if err != nil {
		return err
	}

	if role != nil {
		roles = []types.Role{*role}
	}

	return output.Print(os.Stdout, NewPrinter(roles, roleOptions.LastUsed))
}

func (g *Getter) GetAllRoles(getRoleDetails bool) (roles []types.Role, err error) {
	roles, err = g.iamClient.ListRoles()
	if err != nil {
		return nil, err
	}

	if getRoleDetails {
		return g.getDetailedRoles(roles)
	}

	return roles, nil
}

func (g *Getter) GetIrsaRolesForCluster(cluster *ekstypes.Cluster, getRoleDetails bool) ([]types.Role, error) {
	oidc, err := g.oidcGetter.GetOidcProviderByCluster(cluster)
	if err != nil {
		return []types.Role{}, err
	}

	roles, err := g.GetAllRoles(getRoleDetails)
	if err != nil {
		return []types.Role{}, err
	}

	irsaRoles := []types.Role{}
	providerUrlEscaped := url.QueryEscape(awssdk.ToString(oidc.Url))

	for _, r := range roles {
		if strings.Contains(awssdk.ToString(r.AssumeRolePolicyDocument), providerUrlEscaped) {
			irsaRoles = append(irsaRoles, r)
		}
	}

	return irsaRoles, nil
}

func (g *Getter) GetRoleByName(name string) (*types.Role, error) {
	role, err := g.iamClient.GetRole(name)

	if err != nil {
		var nsee *types.NoSuchEntityException
		if errors.As(err, &nsee) {
			return nil, resource.NotFoundError(fmt.Sprintf("iam-role %q not found", name))
		}
		return nil, err
	}

	return role, nil
}

func (g *Getter) GetRolesByNameSearch(nameSearch string) ([]types.Role, error) {
	roles, err := g.iamClient.ListRoles()
	if err != nil {
		return nil, err
	}

	filtered := []types.Role{}
	n := strings.ToLower(nameSearch)

	for _, r := range roles {
		if strings.Contains(strings.ToLower(awssdk.ToString(r.RoleName)), n) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("no iam-role found searching for %q", nameSearch))
	}

	return filtered, nil
}

func (g *Getter) getDetailedRoles(roles []types.Role) ([]types.Role, error) {
	detailedRoles := make([]types.Role, 0, len(roles))
	for _, r := range roles {
		role, err := g.iamClient.GetRole(awssdk.ToString(r.RoleName))
		if err != nil {
			return []types.Role{}, err
		}

		detailedRoles = append(detailedRoles, *role)
	}

	return detailedRoles, nil
}
