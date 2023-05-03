package iam_policy

import (
	"errors"
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/iam_role"
)

type Getter struct {
	iamClient  *aws.IAMClient
	roleGetter *iam_role.Getter
}

func NewGetter(iamClient *aws.IAMClient) *Getter {
	return &Getter{iamClient, iam_role.NewGetter(iamClient)}
}

func (g *Getter) Init() {
	if g.iamClient == nil {
		g.iamClient = aws.NewIAMClient()
	}
	g.roleGetter = iam_role.NewGetter(g.iamClient)
}

func (g *Getter) Get(arn string, output printer.Output, options resource.Options) error {
	policyOptions, ok := options.(*IamPolicyOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to IamPolicyOptions")
	}

	var policy *types.Policy
	var policies []types.Policy
	var err error

	if policyOptions.Role != "" {
		policies, err := g.GetPoliciesByRoleName(policyOptions.Role)
		if err != nil {
			return err
		}
		return output.Print(os.Stdout, NewRolePolicyPrinter(policies))
	}

	if arn != "" {
		policy, err = g.GetPolicyByArn(arn)
	} else if policyOptions.NameSearch != "" {
		policies, err = g.GetPoliciesByNameSearch(policyOptions.NameSearch)
	} else {
		policies, err = g.GetAllPolicies()
	}

	if err != nil {
		return err
	}

	if policy != nil {
		policies = []types.Policy{*policy}
	}

	return output.Print(os.Stdout, NewPrinter(policies))
}

func (g *Getter) GetAllPolicies() (policies []types.Policy, err error) {
	policies, err = g.iamClient.ListPolicies("")
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (g *Getter) GetPolicyByArn(arn string) (*types.Policy, error) {
	policy, err := g.iamClient.GetPolicy(arn)

	if err != nil {
		var nsee *types.NoSuchEntityException
		if errors.As(err, &nsee) {
			return nil, resource.NotFoundError(fmt.Sprintf("iam-policy %q not found", arn))
		}
		return nil, aws.FormatError(err)
	}

	return policy, nil
}

func (g *Getter) GetPoliciesByNameSearch(nameSearch string) ([]types.Policy, error) {
	policies, err := g.iamClient.ListPolicies("")
	if err != nil {
		return nil, err
	}

	filtered := []types.Policy{}
	n := strings.ToLower(nameSearch)

	for _, p := range policies {
		if strings.Contains(strings.ToLower(awssdk.ToString(p.PolicyName)), n) {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("no iam-policy found searching for %q", nameSearch))
	}

	return filtered, nil
}

func (g *Getter) GetPoliciesByRoleName(roleName string) (*RolePolicies, error) {
	_, err := g.roleGetter.GetRoleByName(roleName)
	if err != nil {
		return nil, err
	}

	// Inline Policies
	inlineNames, err := g.iamClient.ListRolePolicies(roleName)
	if err != nil {
		return nil, err
	}

	inline := make([]InlinePolicy, 0, len(inlineNames))

	for _, policyName := range inlineNames {
		doc, err := g.iamClient.GetRolePolicy(policyName, roleName)
		if err != nil {
			return nil, err
		}

		inline = append(inline, InlinePolicy{policyName, doc})
	}

	// Managed Policies
	attached, err := g.iamClient.ListAttachedRolePolicies(roleName)
	if err != nil {
		return nil, err
	}

	managed := make([]ManagedPolicy, 0, len(attached))

	for _, p := range attached {
		arn := awssdk.ToString(p.PolicyArn)
		policy, err := g.iamClient.GetPolicy(arn)
		if err != nil {
			return nil, err
		}

		doc, err := g.iamClient.GetPolicyVersion(arn, awssdk.ToString(policy.DefaultVersionId))
		if err != nil {
			return nil, err
		}

		managed = append(managed, ManagedPolicy{policy, doc})
	}

	return &RolePolicies{inline, managed}, nil
}
