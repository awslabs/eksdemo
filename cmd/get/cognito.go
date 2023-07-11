package get

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/userpool"
	"github.com/spf13/cobra"
)

var cognitoResources []func() *resource.Resource

func NewGetCognitoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cognito",
		Short: "Amazon Cognito Resources",
	}

	// Don't show flag errors for `get cognito` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range cognitoResources {
		cmd.AddCommand(r().NewGetCmd())
	}

	return cmd
}

func init() {
	cognitoResources = []func() *resource.Resource{
		userpool.New,
	}
}
