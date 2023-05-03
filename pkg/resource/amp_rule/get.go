package amp_rule

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/amp_workspace"
)

type Getter struct {
	prometheusClient *aws.AMPClient
	workspaceGetter  *amp_workspace.Getter
}

func NewGetter(prometheusClient *aws.AMPClient) *Getter {
	return &Getter{prometheusClient, amp_workspace.NewGetter(prometheusClient)}
}

func (g *Getter) Init() {
	if g.prometheusClient == nil {
		g.prometheusClient = aws.NewAMPClient()
	}
	g.workspaceGetter = amp_workspace.NewGetter(g.prometheusClient)
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	ruleOptions, ok := options.(*AmpRuleOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to AmpRuleOptions")
	}

	workspaceId := ruleOptions.WorkspaceId

	if ruleOptions.Alias != "" {
		workspace, err := g.workspaceGetter.GetAmpByAlias(ruleOptions.Alias)
		if err != nil {
			return err
		}
		workspaceId = awssdk.ToString(workspace.Workspace.WorkspaceId)
	}

	var rule *types.RuleGroupsNamespaceDescription
	var rules []*types.RuleGroupsNamespaceDescription
	var err error

	if name == "" {
		rules, err = g.GetAllRules(workspaceId)
	} else {
		rule, err = g.GetRuleByName(name, workspaceId)
		rules = []*types.RuleGroupsNamespaceDescription{rule}
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(rules))
}

func (g *Getter) GetAllRules(workspaceId string) ([]*types.RuleGroupsNamespaceDescription, error) {
	ruleSummaries, err := g.prometheusClient.ListRuleGroupsNamespaces(workspaceId)
	if err != nil {
		return nil, err
	}

	rules := make([]*types.RuleGroupsNamespaceDescription, 0, len(ruleSummaries))

	for _, rs := range ruleSummaries {
		result, err := g.prometheusClient.DescribeRuleGroupsNamespace(awssdk.ToString(rs.Name), workspaceId)
		if err != nil {
			return nil, err
		}
		rules = append(rules, result)
	}

	return rules, nil
}

func (g *Getter) GetRuleByName(ruleName, workspaceId string) (*types.RuleGroupsNamespaceDescription, error) {
	result, err := g.prometheusClient.DescribeRuleGroupsNamespace(ruleName, workspaceId)
	if err != nil {
		return nil, aws.FormatError(err)
	}

	return result, nil
}
