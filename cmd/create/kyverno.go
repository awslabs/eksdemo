package create

import (
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/kyverno"
	"github.com/spf13/cobra"
)

var kyvernoPolicies []func() *resource.Resource

func NewKyvernoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kyverno",
		Short: "Kyverno Policy",
	}

	// Don't show flag errors for `create kyverno` without a subcommand
	cmd.DisableFlagParsing = true

	for _, r := range kyvernoPolicies {
		cmd.AddCommand(r().NewCreateCmd())
	}

	return cmd
}

func init() {
	kyvernoPolicies = []func() *resource.Resource{
		kyverno.NewRequireRequestsPolicy,
		kyverno.NewRestrictAnonymousPolicy,
	}
}
