package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watch-tags",
	Short: "A CLI tool to watch Kubernetes deployment image tags",
	Long:  `watch-tags is a command-line tool to monitor container image tags of deployments in a Kubernetes namespace.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
