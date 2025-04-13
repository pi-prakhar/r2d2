package cmd

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/spf13/cobra"
)

var updateTagCmd = &cobra.Command{
	Use:   "update-tag",
	Short: "Update the image tag of selected services in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(services) == 0 || tag == "" {
			return fmt.Errorf("namespace, services, and tag are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %v", err)
		}

		// Update each service with the new tag
		for _, svc := range services {
			fmt.Printf("Updating service %s in namespace %s to tag %s...\n", svc, namespace, tag)
			err := k8s.UpdateDeploymentTag(clientset, namespace, svc, tag)
			if err != nil {
				fmt.Printf("Failed to update %s: %v\n", svc, err)
			} else {
				fmt.Printf("Updated %s successfully!\n", svc)
			}
		}

		return nil
	},
}

func init() {
	updateTagCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	updateTagCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of service/deployment names (required)")
	updateTagCmd.Flags().StringVarP(&tag, "tag", "t", "", "Image tag to deploy (required)")
	rootCmd.AddCommand(updateTagCmd)
}
