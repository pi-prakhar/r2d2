package cmd

import (
	"fmt"
	"os"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/spf13/cobra"
)

var watchLogsCmd = &cobra.Command{
	Use:   "watch-logs",
	Short: "Watch logs for Kubernetes pods",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 {
			return fmt.Errorf("--namespace and --services are required")
		}

		if path == "" {
			// Use current working directory if --location not provided
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			path = cwd
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		for _, pod := range names {
			go func(p string) {
				err := k8s.GetLogs(clientset, namespace, p, path)
				if err != nil {
					fmt.Printf("Log error for %s: %v\n", p, err)
				}
			}(pod)
		}
		select {}
	},
}

func init() {
	watchLogsCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	watchLogsCmd.Flags().StringSliceVarP(&names, "names", "p", []string{}, constants.CommonFlagDescPodNames)
	watchLogsCmd.Flags().StringVarP(&path, "location", "l", "", constants.CommonFlagDescLocation)
	rootCmd.AddCommand(watchLogsCmd)
}
