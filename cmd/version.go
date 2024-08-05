package cmd

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and build information for eksdemo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("eksdemo: %#v\n", version.GetVersionInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
