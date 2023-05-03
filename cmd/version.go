package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

type Version struct {
	Version string
	Date    string
	Commit  string
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version of eksdemo",
	Run: func(cmd *cobra.Command, args []string) {
		eksdemoVersion := Version{
			Version: version,
			Date:    date,
			Commit:  commit,
		}

		fmt.Printf("eksdemo version info: %#v\n", eksdemoVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
