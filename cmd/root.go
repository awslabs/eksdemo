package cmd

import (
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/cmd/create"
	del "github.com/awslabs/eksdemo/cmd/delete"
	"github.com/awslabs/eksdemo/cmd/get"
	"github.com/awslabs/eksdemo/cmd/install"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/helm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile, region, profile string
var debug, responseBodyDebug bool

var rootCmd = &cobra.Command{
	Use:              "eksdemo",
	Short:            "The easy button for learning, testing, and demoing Amazon EKS",
	PersistentPreRun: preRun,
	SilenceErrors:    true,
	Long: `The easy button for learning, testing, and demoing Amazon EKS:
  * Install complex applications and dependencies with a single command
  * Extensive application catalog (over 50 CNCF, open source and related projects)
  * Customize application installs easily with simple command line flags
  * Query and search AWS resources with over 60 kubectl-like get commands`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func preRun(cmd *cobra.Command, args []string) {
	// This will work in the future if the issue below is fixed:
	// https://github.com/spf13/cobra/issues/1413
	// cmd.SilenceUsage = true

	aws.Init(profile, region, debug, responseBodyDebug)
	helm.Init(debug)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		newCmdCompletion(rootCmd),
		create.NewCreateCmd(),
		del.NewDeleteCmd(),
		get.NewGetCmd(),
		install.NewInstallCmd(),
		install.NewUninstallCmd(),
		newCmdUpdate(),
	)

	// TODO: implement configuration
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.eksdemo.yaml)")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "use the specific profile from your credential file")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "the region to use, overrides config/env settings")

	// Hidden debug flags
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", debug, "")
	rootCmd.PersistentFlags().BoolVar(&responseBodyDebug, "debug-response", responseBodyDebug, "")
	rootCmd.PersistentFlags().MarkHidden("debug")
	rootCmd.PersistentFlags().MarkHidden("debug-response")

	// Hide help command
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".eksdemo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".eksdemo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
