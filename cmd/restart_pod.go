package cmd

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/internal/k8s"
	"github.com/spf13/cobra"
)

var restartPodCmd = &cobra.Command{
	Use:   "pod",
	Short: "Restart pods (by deleting them; they'll auto-restart if part of a deployment/replica set)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 {
			return fmt.Errorf("both --namespace and --names are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("failed to get Kubernetes client: %v", err)
		}

		for _, pod := range names {
			err := k8s.RestartPod(clientset, namespace, pod)
			if err != nil {
				fmt.Printf("Failed to restart pod %s: %v\n", pod, err)
			} else {
				fmt.Printf("Restarted pod %s\n", pod)
			}
		}
		return nil
	},
}

func init() {
	restartPodCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	restartPodCmd.Flags().StringSliceVarP(&names, "names", "p", []string{}, constants.CommonFlagDescPodNames)
	restartPodCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	restartPodCmd.RegisterFlagCompletionFunc("names", getDeployments)
}
