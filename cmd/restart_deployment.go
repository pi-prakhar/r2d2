package cmd

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/spf13/cobra"
)

var restartDeploymentCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Restart a deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 {
			return fmt.Errorf("both --namespace and --names are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("failed to get Kubernetes client: %v", err)
		}

		for _, deployment := range names {
			err := k8s.RestartDeployment(clientset, namespace, deployment)
			if err != nil {
				fmt.Printf("Failed to restart deployment %s: %v\n", deployment, err)
			} else {
				fmt.Printf("Restarted deployment %s\n", deployment)
			}
		}
		return nil
	},
}

func init() {
	restartDeploymentCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	restartDeploymentCmd.Flags().StringSliceVarP(&names, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	restartDeploymentCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	restartDeploymentCmd.RegisterFlagCompletionFunc("names", getDeployments)
}
