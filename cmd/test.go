package cmd

import (
	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:    "test",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {

		cmd.SilenceUsage = true

		return nil

	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
