package cmd

import (
	"fmt"
	"time"

	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/internal/k8s"
	"github.com/pi-prakhar/r2d2/ui/table"
	"github.com/spf13/cobra"
)

var watchImagesCmd = &cobra.Command{
	Use:   "watch-images",
	Short: "Watches deployment images in a namespace",
	RunE: func(cmd *cobra.Command, args []string) error {
		if namespace == "" || len(names) == 0 {
			return fmt.Errorf("--namespace and --names are required")
		}

		clientset, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		app := table.NewWatchImagesApp(namespace)

		go func() {
			for {
				data, err := k8s.FetchDeploymentInfo(clientset, namespace, names)
				if err != nil {
					app.Stop()
					fmt.Printf("error fetching deployment tags: %v\n", err)
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

	watchImagesCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	watchImagesCmd.Flags().StringSliceVarP(&names, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	watchImagesCmd.Flags().IntVarP(&frequency, "frequency", "f", constants.DeploymentWatchImagesDefaultFrequency,
		fmt.Sprintf(constants.CommonFlagDescWatchFrequency, "images"))
	watchImagesCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	watchImagesCmd.RegisterFlagCompletionFunc("names", getDeployments)
	rootCmd.AddCommand(watchImagesCmd)
}
