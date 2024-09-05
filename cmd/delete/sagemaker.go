package delete

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/sagemaker/userprofile"
	"github.com/spf13/cobra"
)

var sagemaker []func() *resource.Resource

func NewSageMakerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sagemaker",
		Short:   "Amazon SageMaker Resources",
		Aliases: []string{"sm"},
	}

	// Don't show flag errors for `delete sagemaker` without a subcommand
	cmd.DisableFlagParsing = true

	for _, sm := range sagemaker {
		cmd.AddCommand(sm().NewDeleteCmd())
	}

	return cmd
}

func init() {
	sagemaker = []func() *resource.Resource{
		userprofile.New,
	}
}
