package install

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/application/crossplane/core"
	"github.com/awslabs/eksdemo/pkg/application/crossplane/ec2"
	"github.com/awslabs/eksdemo/pkg/application/crossplane/iam"
	"github.com/awslabs/eksdemo/pkg/application/crossplane/s3"
	"github.com/spf13/cobra"
)

var crossplane []func() *application.Application

func NewInstallCrossplaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crossplane",
		Short:   "The Cloud Native Control Plane",
		Aliases: []string{"cp"},
	}

	// Don't show flag errors for `install crossplane` without a subcommand
	cmd.DisableFlagParsing = true

	for _, cp := range crossplane {
		cmd.AddCommand(cp().NewInstallCmd())
	}

	return cmd
}

func NewUninstallCrossplaneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "crossplane",
		Short:   "The Cloud Native Control Plane",
		Aliases: []string{"cp"},
	}

	// Don't show flag errors for `uninstall crossplane` without a subcommand
	cmd.DisableFlagParsing = true

	for _, cp := range crossplane {
		cmd.AddCommand(cp().NewUninstallCmd())
	}

	return cmd
}

func init() {
	crossplane = []func() *application.Application{
		core.NewApp,
		ec2.NewApp,
		iam.NewApp,
		s3.NewApp,
	}
}
