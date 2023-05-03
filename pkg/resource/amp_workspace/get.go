package amp_workspace

import (
	"errors"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type AmpWorkspace struct {
	Workspace        *types.WorkspaceDescription
	WorkspaceLogging *types.LoggingConfigurationMetadata
}

type Getter struct {
	prometheusClient *aws.AMPClient
}

func NewGetter(prometheusClient *aws.AMPClient) *Getter {
	return &Getter{prometheusClient}
}

func (g *Getter) Init() {
	if g.prometheusClient == nil {
		g.prometheusClient = aws.NewAMPClient()
	}
}

func (g *Getter) Get(alias string, output printer.Output, options resource.Options) error {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(workspaces))
}

func (g *Getter) GetAll(alias string) ([]AmpWorkspace, error) {
	ampSummaries, err := g.prometheusClient.ListWorkspaces(alias)
	if err != nil {
		return nil, err
	}

	workspaces := make([]AmpWorkspace, 0, len(ampSummaries))

	for _, summary := range ampSummaries {
		// ListWorkspaces API will return workspaces that begin with alias, so drop those that don't match exactly
		if alias != "" && awssdk.ToString(summary.Alias) != alias {
			continue
		}

		workspace, err := g.prometheusClient.DescribeWorkspace(awssdk.ToString(summary.WorkspaceId))
		if err != nil {
			return nil, err
		}

		logging, err := g.prometheusClient.DescribeLoggingConfiguration(awssdk.ToString(summary.WorkspaceId))
		// Return all errors except NotFound
		var rnfe *types.ResourceNotFoundException
		if err != nil && !errors.As(err, &rnfe) {
			return nil, err
		}

		workspaces = append(workspaces, AmpWorkspace{workspace, logging})
	}

	if alias != "" && len(workspaces) == 0 {
		return nil, resource.NotFoundError(fmt.Sprintf("workspace alias %q not found", alias))
	}

	return workspaces, nil
}

func (g *Getter) GetAmpByAlias(alias string) (AmpWorkspace, error) {
	workspaces, err := g.GetAll(alias)
	if err != nil {
		return AmpWorkspace{}, err
	}

	found := []AmpWorkspace{}

	for _, w := range workspaces {
		if w.Workspace.Status.StatusCode != types.WorkspaceStatusCodeDeleting {
			found = append(found, w)
		}
	}

	if len(found) == 0 {
		return AmpWorkspace{}, resource.NotFoundError(fmt.Sprintf("workspace alias %q not found", alias))
	}

	if len(found) > 1 {
		return AmpWorkspace{}, fmt.Errorf("multiple workspaces found with alias: %s", alias)
	}

	return found[0], nil
}
