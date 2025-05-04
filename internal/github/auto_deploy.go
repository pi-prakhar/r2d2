package ghservice

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/pi-prakhar/r2d2/config"
	"github.com/pi-prakhar/r2d2/constants"
	"github.com/pi-prakhar/r2d2/utils/helper"
)

func GetGithubWorkflowStatus(ghClient *github.Client, config config.AutoDeployConfig) (constants.AutoDeployStatus, error) {
	ctx := context.Background()

	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	runs, _, err := ghClient.Actions.ListRepositoryWorkflowRuns(ctx, config.Owner, config.Repo, opts)
	if err != nil {
		return constants.AutoDeployStatusFailed, err
	}

	filteredRuns := []*github.WorkflowRun{}
	for _, run := range runs.WorkflowRuns {
		if (run.HeadBranch != nil && *run.HeadBranch == config.Tag) ||
			(run.HeadSHA != nil && *run.HeadSHA == config.Tag) {
			filteredRuns = append(filteredRuns, run)
		}
	}

	if len(filteredRuns) == 0 {
		return constants.AutoDeployStatusWaiting, nil
	}

	var ecrPushSuccessful bool

	for _, run := range filteredRuns {
		if run.Name != nil && strings.Contains(strings.ToLower(*run.Name), "ecr") {
			if *run.Status == "completed" {
				if *run.Conclusion == "success" {
					ecrPushSuccessful = true
				}
			} else {
				return constants.AutoDeployStatusInProgress, nil
			}
			break
		}
	}

	if !ecrPushSuccessful {
		return constants.AutoDeployStatusFailed, fmt.Errorf("ECR Push workflow failed. Aborting deployments")
	}

	return constants.AutoDeployStatusCompleted, nil
}
func ParseAutoDeployFlags(
	githubRepository string,
	tag string,
	namespace string,
	names []string,
	frequency int,
	owner string,
) (config.AutoDeployConfig, error) {
	token, err := helper.GetGitHubTokenFromEnv()
	if err != nil {
		return config.AutoDeployConfig{}, fmt.Errorf("failed to get GitHub token: %w", err)
	}
	return config.AutoDeployConfig{
		Owner:       owner,
		Repo:        githubRepository,
		Tag:         tag,
		Namespace:   namespace,
		Deployments: names,
		Interval:    frequency,
		Token:       token,
	}, nil
}
