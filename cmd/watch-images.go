package cmd

import (
	"fmt"
	"time"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/spf13/cobra"
)

var watchImagesCmd = &cobra.Command{
	Use:   "watch-images",
	Short: "Watches deployment images in a namespace",
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

		app := utils.NewWatchImagesApp(namespace)

		go func() {
			for {
				data, err := k8s.FetchDeploymentInfo(clientset, namespace, services)
				if err != nil {
					app.Stop()
					fmt.Printf("Error fetching deployment tags: %v\n", err)
					return
				}

				app.UpdateTable(data)
				time.Sleep(time.Duration(frequency) * time.Second)
			}
		}()

		if err := app.Run(); err != nil {
			fmt.Printf("Error running the application: %v\n", err)
		}
	},
}

func init() {
	watchImagesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	watchImagesCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of service/deployment names (required)")
	watchImagesCmd.Flags().IntVarP(&frequency, "frequency", "f", 60, "Frequency of fetching images in seconds (default: 60)")
	rootCmd.AddCommand(watchImagesCmd)
}
