package cmd

import (
	"fmt"
	"time"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/spf13/cobra"
)

var (
	namespace string
	services  []string
)

var watchTagsCmd = &cobra.Command{
	Use:   "watch-tags",
	Short: "Watches deployment tags in a namespace",
	Run: func(cmd *cobra.Command, args []string) {
		if namespace == "" || len(services) == 0 {
			fmt.Println("Error: --namespace and --services are required")
			return
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			fmt.Printf("Error creating Kubernetes client: %v\n", err)
			return
		}

		app := utils.NewWatchTagsApp()

		go func() {
			for {
				data, err := k8s.FetchDeploymentTags(clientset, namespace, services)
				if err != nil {
					app.Stop()
					fmt.Printf("Error fetching deployment tags: %v\n", err)
					return
				}

				app.UpdateTable(data)
				time.Sleep(60 * time.Second)
			}
		}()

		if err := app.Run(); err != nil {
			fmt.Printf("Error running the application: %v\n", err)
		}
	},
}

func init() {
	watchTagsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	watchTagsCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of service/deployment names (required)")
	rootCmd.AddCommand(watchTagsCmd)
}
