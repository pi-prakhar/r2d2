package cmd

import (
	"fmt"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/internal/k8s"
	"github.com/spf13/cobra"
)

var updateTagCmd = &cobra.Command{
	Use:   "update-tag",
	Short: "Update the image tag of selected deployments in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 || tag == "" {
			return fmt.Errorf("--namespace, --names, --tag are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %v", err)
		}

		// Update each deployment with the new tag
		for _, svc := range names {
			fmt.Printf("Updating deployment %s in namespace %s to tag %s...\n", svc, namespace, tag)
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
	updateTagCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	updateTagCmd.Flags().StringSliceVarP(&names, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	updateTagCmd.Flags().StringVarP(&tag, "tag", "t", "", constants.CommonFlagDescTag)
	updateTagCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	updateTagCmd.RegisterFlagCompletionFunc("names", getDeployments)
	rootCmd.AddCommand(updateTagCmd)
}
