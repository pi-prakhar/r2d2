package cmd

import (
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart Kubernetes resources like deployments or pods",
}

func init() {
	rootCmd.AddCommand(restartCmd)
	restartCmd.AddCommand(restartDeploymentCmd)
	restartCmd.AddCommand(restartPodCmd)
}
