package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/google/go-github/v45/github"
	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/k8s"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"k8s.io/client-go/kubernetes"
)

type AutoDeployConfig struct {
	owner       string
	repo        string
	tag         string
	namespace   string
	deployments []string
	interval    int
	token       string
}

var (
	autoDeployNamespace string
	autoDeployNames     []string
	autoDeployTag       string
	autoDeployInterval  int
	autoDeployRepo      string
)

var autoDeployCmd = &cobra.Command{
	Use:   "auto-deploy",
	Short: "Automatically deploy GitHub releases to Kubernetes",
	Long: `Automatically watches GitHub workflow status for a specific tag and updates Kubernetes deployments
when the workflows complete successfully.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := parseAutoDeployFlags()

		// Create GitHub client
		ghClient := createGitHubClient(config.token)

		// Create Kubernetes client
		k8sClient, err := k8s.GetClientSet()
		if err != nil {
			return fmt.Errorf("error creating Kubernetes client: %w", err)
		}

		// Start watching
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = fmt.Sprintf(" Watching GitHub workflows for tag: %s...", config.tag)
		s.Start()

		for {
			completed := checkAndUpdateDeployments(ghClient, k8sClient, config, config.tag)
			if completed {
				s.Stop()
				color.Green("\nAll workflows completed and deployments updated!")
				return nil
			}

			time.Sleep(time.Duration(config.interval) * time.Second)
		}
	},
}

func parseAutoDeployFlags() AutoDeployConfig {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		color.Red("Error: GITHUB_TOKEN environment variable is not set")
		os.Exit(1)
	}

	return AutoDeployConfig{
		owner:       "Orange-Health",
		repo:        autoDeployRepo,
		tag:         autoDeployTag,
		namespace:   autoDeployNamespace,
		deployments: autoDeployNames,
		interval:    autoDeployInterval,
		token:       token,
	}
}

func createGitHubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func checkAndUpdateDeployments(ghClient *github.Client, k8sClient *kubernetes.Clientset, config AutoDeployConfig, tag string) bool {
	ctx := context.Background()

	// Check GitHub workflow status
	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	runs, _, err := ghClient.Actions.ListRepositoryWorkflowRuns(ctx, config.owner, config.repo, opts)
	if err != nil {
		color.Red("⚠️  Error fetching workflow runs: %v", err)
		return false
	}

	// Filter runs for the specific tag
	filteredRuns := []*github.WorkflowRun{}
	for _, run := range runs.WorkflowRuns {
		if (run.HeadBranch != nil && *run.HeadBranch == tag) ||
			(run.HeadSHA != nil && *run.HeadSHA == tag) {
			filteredRuns = append(filteredRuns, run)
		}
	}

	if len(filteredRuns) == 0 {
		color.Yellow("⏳ No workflow runs found for tag '%s'. They might not have started yet.", tag)
		return false
	}

	// Clear screen and display header
	fmt.Print("\033[H\033[2J") // Clear screen
	bold := color.New(color.Bold).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()

	// Display header with tag and repository info
	fmt.Printf("\n🚀 %s %s\n", bold("Deployment Monitor:"), cyan(tag))
	fmt.Printf("📦 %s/%s\n", config.owner, config.repo)
	fmt.Println(dim(strings.Repeat("─", 80)))

	// Display workflow status
	allCompleted := true
	totalRuns := len(filteredRuns)
	completedRuns := 0
	successfulRuns := 0

	for _, run := range filteredRuns {
		var statusColor *color.Color
		var statusEmoji string
		var statusText string

		completed := *run.Status == "completed"
		success := completed && *run.Conclusion == "success"

		if completed {
			completedRuns++
			if success {
				successfulRuns++
				statusColor = color.New(color.FgGreen)
				statusEmoji = "✅"
				statusText = "Success"
			} else {
				statusColor = color.New(color.FgRed)
				statusEmoji = "❌"
				statusText = *run.Conclusion
			}
		} else {
			statusColor = color.New(color.FgYellow)
			statusEmoji = "🔄"
			statusText = "In Progress"
			allCompleted = false
		}

		workflowName := *run.Name
		createdAt := run.CreatedAt.Time.Format("15:04:05")

		// Display workflow status with emoji and better formatting
		fmt.Printf("\n%s %s\n", statusEmoji, bold(workflowName))
		fmt.Printf("   %s %s\n", statusColor.Sprint("Status:"), statusColor.Sprint(statusText))
		fmt.Printf("   %s %s\n", dim("Started:"), createdAt)
	}

	// Display progress summary
	fmt.Println(dim(strings.Repeat("─", 80)))
	fmt.Printf("Progress: %d/%d completed", completedRuns, totalRuns)
	if completedRuns > 0 {
		fmt.Printf(" (%d successful)\n", successfulRuns)
	} else {
		fmt.Println()
	}

	// Display GitHub Actions URL
	fmt.Printf("\n🔗 %s\n", color.BlueString(fmt.Sprintf("https://github.com/%s/%s/actions", config.owner, config.repo)))

	// If all workflows are completed successfully, update Kubernetes deployments
	if allCompleted {
		fmt.Println(dim(strings.Repeat("─", 80)))
		color.Green("\n🎉 All workflows completed successfully!")
		color.Cyan("\n📡 Updating Kubernetes deployments...")

		successCount := 0
		for _, deployment := range config.deployments {
			err := k8s.UpdateDeploymentTag(k8sClient, config.namespace, deployment, tag)
			if err != nil {
				color.Red("❌ Failed to update %s: %v", deployment, err)
				continue
			}
			color.Green("✅ Updated deployment: %s", deployment)
			successCount++
		}

		// Display final summary
		fmt.Println(dim(strings.Repeat("─", 80)))
		if successCount == len(config.deployments) {
			color.Green("✨ All deployments successfully updated to %s!", tag)
		} else {
			color.Yellow("⚠️  %d/%d deployments updated to %s", successCount, len(config.deployments), tag)
		}
	}

	return allCompleted
}

func init() {
	autoDeployCmd.Flags().StringVarP(&autoDeployNamespace, "namespace", "n", "", constants.CommonFlagDescNamespace)
	autoDeployCmd.Flags().StringSliceVarP(&autoDeployNames, "names", "d", []string{}, constants.CommonFlagDescDeploymentNames)
	autoDeployCmd.Flags().StringVarP(&autoDeployTag, "tag", "t", "", "Git tag to watch (required)")
	autoDeployCmd.Flags().IntVarP(&autoDeployInterval, "interval", "i", 10, "Polling interval in seconds")
	autoDeployCmd.Flags().StringVarP(&autoDeployRepo, "repo", "r", "", "GitHub repository name (required)")

	autoDeployCmd.MarkFlagRequired("repo")
	autoDeployCmd.MarkFlagRequired("tag")

	autoDeployCmd.RegisterFlagCompletionFunc("namespace", getNamespaces)
	autoDeployCmd.RegisterFlagCompletionFunc("names", getDeployments)

	rootCmd.AddCommand(autoDeployCmd)
}
