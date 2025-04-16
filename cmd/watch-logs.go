package cmd

import (
	"fmt"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/spf13/cobra"
	"time"
)

var watchLogsCmd = &cobra.Command{
	Use:   "log",
	Short: "gets logs of a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(services) == 0 || path == "" {
			return fmt.Errorf("--namespace and --services are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		app := utils.NewWatchLogsApp(namespace)

		go func() {
			for {
				data := []k8s.Info{}
				for _, service := range services {
					filePath, err := k8s.GetLogs(clientset, namespace, service, path)
					if err != nil {
						app.Stop()
						fmt.Printf("error fetching pod logs: %v\n", err)
						return
					}
					data = append(data, k8s.Info{
						PodName: service,
						Path:    filePath,
					})

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
	watchLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	watchLogsCmd.Flags().StringSliceVarP(&services, "pods", "p", []string{}, "List of service/deployment names (required)")
	watchLogsCmd.Flags().IntVarP(&frequency, "frequency", "f", 60, "Frequency of fetching logs in seconds (default: 60)")
	watchLogsCmd.Flags().StringVarP(&path, "location", "l", "", "path/location to the file (required)")
	rootCmd.AddCommand(watchLogsCmd)
}
