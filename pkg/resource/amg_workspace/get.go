package amg_workspace

import (
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/grafana/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	grafanaClient *aws.GrafanaClient
}

func NewGetter(grafanaClient *aws.GrafanaClient) *Getter {
	return &Getter{grafanaClient}
}

func (g *Getter) Init() {
	if g.grafanaClient == nil {
		g.grafanaClient = aws.NewGrafanaClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var workspaces []*types.WorkspaceDescription
	var err error

	if name == "" {
		workspaces, err = g.GetAll()
	} else {
		workspaces, err = g.GetAllAmgByName(name)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAll() ([]*types.WorkspaceDescription, error) {
	amgSummaries, err := g.grafanaClient.ListWorkspaces()
	workspaces := make([]*types.WorkspaceDescription, 0, len(amgSummaries))

	if err != nil {
		return nil, err
	}

	for _, summary := range amgSummaries {
		result, err := g.grafanaClient.DescribeWorkspace(awssdk.ToString(summary.Id))
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	return workspaces, nil
}

func (g *Getter) GetAllAmgByName(name string) ([]*types.WorkspaceDescription, error) {
	summaries, err := g.grafanaClient.ListWorkspaces()
	if err != nil {
		return nil, err
	}

	workspaces := make([]*types.WorkspaceDescription, 0, len(summaries))

	for _, s := range summaries {
		if awssdk.ToString(s.Name) != name {
			continue
		}

		result, err := g.grafanaClient.DescribeWorkspace(awssdk.ToString(s.Id))
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, result)
	}

	if len(workspaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace name %q not found", name))
	}

	return workspaces, nil
}

func (g *Getter) GetAmgByName(name string) (*types.WorkspaceDescription, error) {
	workspaces, err := g.GetAllAmgByName(name)
	if err != nil {
		return nil, err
	}

	found := []*types.WorkspaceDescription{}

	for _, w := range workspaces {
		if w.Status != types.WorkspaceStatusDeleting {
			found = append(found, w)
		}
	}

	if len(found) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace name %q not found", name))
	}

	if len(found) > 1 {
		return nil, fmt.Errorf("multiple workspaces found with name: %s", name)
	}

	return found[0], nil
}
