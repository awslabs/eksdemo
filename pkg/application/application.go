package application

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Application struct {
	cmd.Command
	cmd.Flags
	Options

	Dependencies []*resource.Resource
	Installer
	PostInstallResources []*resource.Resource
}

func (a *Application) NewInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     a.Name,
		Short:   a.Description,
		Long:    "Install " + a.Name,
		Aliases: a.Aliases,
		Hidden:  a.Hidden,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.ValidateFlags(cmd, args); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			if a.Common().DryRun {
				a.SetDryRun()
			}

			if err := a.PreDependencies(Install); err != nil {
				return err
			}

			if err := a.CreateDependencies(); err != nil {
				return err
			}

			if err := a.PreInstall(); err != nil {
				return err
			}

			if err := a.Install(a.Options); err != nil {
				return err
			}

			return a.PostInstall(a.Name, a.PostInstallResources)
		},
	}
	a.Flags = a.Options.AddInstallFlags(cmd, a.Flags, a.Installer.Type())

	return cmd
}

func (a *Application) NewUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     a.Name,
		Short:   a.Description,
		Long:    "Uninstall " + a.Name,
		Aliases: a.Aliases,
		Hidden:  a.Hidden,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := a.ValidateFlags(cmd, args); err != nil {
				return err
			}
			cmd.SilenceUsage = true

			if err := a.PreDependencies(Uninstall); err != nil {
				return err
			}

			if a.Options.Common().DeleteDependencies {
				if err := a.DeleteDependencies(); err != nil {
					return err
				}
			}

			return a.Uninstall(a.Options)
		},
	}
	a.Flags = a.Options.AddUninstallFlags(cmd, a.Flags, a.Dependencies != nil)

	return cmd
}
