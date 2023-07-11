package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/userpool"
	"github.com/spf13/cobra"
)

var cognitoResources []func() *resource.Resource

func NewCognitoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cognito",
		Short: "Amazon Cognito Resources",
	}

	// Don't show flag errors for `create cognito` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range cognitoResources {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	cognitoResources = []func() *resource.Resource{
		userpool.New,
	}
}
