package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/pi-prakhar/r2d2/utils"
	"github.com/spf13/cobra"
)

var (
	namespace string
	services  []string
)

var watchCmd = &cobra.Command{
	Use:   "watch-tags",
	Short: "Watches deployment image tags in a namespace",
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

		for {
			utils.ClearTerminal()
			fmt.Printf("Watching image tags for services: %s in namespace: %s\n\n", strings.Join(services, ", "), namespace)

			data, err := k8s.FetchDeploymentTags(clientset, namespace, services)
			if err != nil {
				fmt.Printf("Error fetching deployment tags: %v\n", err)
			} else {
				utils.PrintTable(data)
			}

			time.Sleep(60 * time.Second)
		}
	},
}

func init() {
	watchCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Kubernetes namespace (required)")
	watchCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "List of service/deployment names (required)")
	rootCmd.AddCommand(watchCmd)
}
