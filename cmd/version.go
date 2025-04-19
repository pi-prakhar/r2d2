package cmd

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of your CLI tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:" + version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
