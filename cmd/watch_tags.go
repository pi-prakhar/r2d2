package cmd

import (
	"fmt"
	"time"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils/table"
	"github.com/spf13/cobra"
)

var watchTagsCmd = &cobra.Command{
	Use:   "watch-tags",
	Short: "Watches deployment tags in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 {
			return fmt.Errorf("--namespace, --names are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		if podLevel {
			// Pod level view
			app := table.NewWatchPodTagsApp(namespace)

			go func() {
				for {
					data, err := k8s.FetchPodInfoForDeployments(clientset, namespace, names)
					if err != nil {
						app.Stop()
						fmt.Printf("error fetching pod info: %v\n", err)
						return
					}

					app.UpdateTable(data)
					time.Sleep(time.Duration(frequency) * time.Second)
				}
			}()

			if err := app.Run(); err != nil {
				return fmt.Errorf("error running the application: %w", err)
			}
		} else {
			// Deployment level view (original behavior)
			app := table.NewWatchTagsApp(namespace)

			go func() {
				for {
					data, err := k8s.FetchDeploymentInfo(clientset, namespace, names)
					if err != nil {
						app.Stop()
						fmt.Printf("error fetching deployment info: %v\n", err)
						return
					}

					app.UpdateTable(data)
					time.Sleep(time.Duration(frequency) * time.Second)
				}
			}()

			if err := app.Run(); err != nil {
				return fmt.Errorf("error running the application: %w", err)
			}
		}

		return nil
	},
}

func init() {
	watchTagsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	watchTagsCmd.Flags().StringSliceVarP(&names, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	watchTagsCmd.Flags().IntVarP(&frequency, "frequency", "f", constants.DeploymentWatchTagsDefaultFrequency,
		fmt.Sprintf(constants.CommonFlagDescWatchFrequency, "tags"))
	watchTagsCmd.Flags().BoolVarP(&podLevel, "pod-level", "p", false, constants.CommonFlagDescPodLevel)
	watchTagsCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	watchTagsCmd.RegisterFlagCompletionFunc("names", getDeployments)
	rootCmd.AddCommand(watchTagsCmd)
}
