package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pi-prakhar/r2d2/constants"
	ghservice "github.com/pi-prakhar/r2d2/internal/github"
	"github.com/pi-prakhar/r2d2/internal/k8s"
	"github.com/pi-prakhar/r2d2/ui/view"
	"github.com/pi-prakhar/r2d2/utils/helper"
	"github.com/spf13/cobra"
)

var autoDeployCmd = &cobra.Command{
	Use:   "auto-deploy",
	Short: "Automatically deploy GitHub releases to Kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := ghservice.ParseAutoDeployFlags(
			githubRepository,
			tag,
			namespace,
			names,
			frequency,
			"Orange-Health", // Example owner
		)
		if err != nil {
			return fmt.Errorf("error parsing flags: %w", err)
		}

		ghClient := ghservice.GetSharedGitHubClient(config.Token)
		k8sClient, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		var errors []error
		var summary string

		status, _ := ghservice.GetGithubWorkflowStatus(ghClient, config)

		if status == constants.AutoDeployStatusInProgress || status == constants.AutoDeployStatusWaiting {
			ui := view.NewUI(config.Owner, config.Repo, config.Tag)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			wg := sync.WaitGroup{}
			wg.Add(1)

			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						status, _ := ghservice.GetGithubWorkflowStatus(ghClient, config)

						if status == constants.AutoDeployStatusInProgress || status == constants.AutoDeployStatusWaiting {
							ui.Clear()
							ui.PrintStatus(status == constants.AutoDeployStatusWaiting)
							time.Sleep(time.Duration(config.Interval) * time.Second)
							continue
						}
						ui.Stop()
					}
				}
			}()
			if err := ui.Run(); err != nil {
				fmt.Printf("error running UI: %v\n", err)
				cancel()
				wg.Wait()
				return fmt.Errorf("ui error: %w", err)
			}
			cancel()
			wg.Wait()
		}

		// Final deployment update check
		status, err = ghservice.GetGithubWorkflowStatus(ghClient, config)
		if status == constants.AutoDeployStatusFailed {
			errors = append(errors, err)
			summary = helper.BuildFinalSummary(false, false, config, errors)
		}

		if status == constants.AutoDeployStatusCompleted {
			k8sSuccess, deployErrors := k8s.UpdateK8sDeploymentTag(k8sClient, config)
			if k8sSuccess {
				summary = helper.BuildFinalSummary(true, true, config, deployErrors)
			} else {
				summary = helper.BuildFinalSummary(true, false, config, deployErrors)
			}
		}

		fmt.Println(summary)
		return nil
	},
}

func init() {
	autoDeployCmd.Flags().StringVarP(&namespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	autoDeployCmd.Flags().StringSliceVarP(&names, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	autoDeployCmd.Flags().StringVarP(&tag, "tag", "t", "", constants.AutoDeployGitTagToWatch)
	autoDeployCmd.Flags().IntVarP(&frequency, "frequency", "f", constants.AutoDeployWatchTagDefaultFrequency, fmt.Sprintf(constants.CommonFlagDescWatchFrequency, "GitHub workflows"))
	autoDeployCmd.Flags().StringVarP(&githubRepository, "repository", "r", "", constants.AutoDeployGithubRepository)

	autoDeployCmd.MarkFlagRequired("repository")
	autoDeployCmd.MarkFlagRequired("tag")

	autoDeployCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	autoDeployCmd.RegisterFlagCompletionFunc("names", getDeployments)

	rootCmd.AddCommand(autoDeployCmd)
}
