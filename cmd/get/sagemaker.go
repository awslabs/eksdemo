package get

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/sagemaker/domain"
	"github.com/awslabs/eksdemo/pkg/resource/sagemaker/userprofile"
	"github.com/spf13/cobra"
)

var sagemaker []func() *resource.Resource

func NewGetSageMakerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sagemaker",
		Short:   "Amazon SageMaker Resources",
		Aliases: []string{"sm"},
	}

	// Don't show flag errors for `get sagemaker` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range sagemaker {
		cmd.AddCommand(r().NewGetCmd())
	}

	return cmd
}

func init() {
	sagemaker = []func() *resource.Resource{
		domain.New,
		userprofile.New,
	}
}
