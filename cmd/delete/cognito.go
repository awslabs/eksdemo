package delete

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/client"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/domain"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/userpool"
	"github.com/spf13/cobra"
)

var cognitoResources []func() *resource.Resource

func NewCognitoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cognito",
		Short: "Amazon Cognito Resources",
	}

	// Don't show flag errors for `delete cognito` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range cognitoResources {
		cmd.AddCommand(r().NewDeleteCmd())
	}

	return cmd
}

func init() {
	cognitoResources = []func() *resource.Resource{
		client.New,
		domain.New,
		userpool.New,
	}
}
