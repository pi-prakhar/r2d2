package cmd

import (
	"fmt"
	"time"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/spf13/cobra"
)

var watchTagsCmd = &cobra.Command{
	Use:   "watch-tags",
	Short: "Watches deployment tags in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(services) == 0 {
			return fmt.Errorf("--namespace and --services are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		app := utils.NewWatchTagsApp(namespace)

		go func() {
			for {
				data, err := k8s.FetchDeploymentInfo(clientset, namespace, services)
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

		return nil
	},
}

func init() {
	watchTagsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	watchTagsCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	watchTagsCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of service/deployment names (required)")
	watchTagsCmd.RegisterFlagCompletionFunc("services", getDeployments)
	watchTagsCmd.Flags().IntVarP(&frequency, "frequency", "f", 60, "Frequency of fetching tags in seconds (default: 60)")
	rootCmd.AddCommand(watchTagsCmd)
}
