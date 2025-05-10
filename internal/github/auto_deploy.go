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

type WorkflowJobStatus struct {
	Name       string
	Status     string
	Conclusion string
}

func GetGithubWorkflowStatus(ghClient *github.Client, config config.AutoDeployConfig) (constants.AutoDeployStatus, []WorkflowJobStatus, error) {
	ctx := context.Background()

	opts := &github.ListWorkflowRunsOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	runs, _, err := ghClient.Actions.ListRepositoryWorkflowRuns(ctx, config.Owner, config.Repo, opts)
	if err != nil {
		return constants.AutoDeployStatusFailed, nil, err
	}

	filteredRuns := []*github.WorkflowRun{}
	for _, run := range runs.WorkflowRuns {
		if (run.HeadBranch != nil && *run.HeadBranch == config.Tag) ||
			(run.HeadSHA != nil && *run.HeadSHA == config.Tag) {
			filteredRuns = append(filteredRuns, run)
		}
	}

	if len(filteredRuns) == 0 {
		return constants.AutoDeployStatusWaiting, nil, nil
	}

	var ecrPushSuccessful bool

	for _, run := range filteredRuns {
		if strings.Contains(strings.ToLower(run.GetName()), "ecr") {
			if run.GetStatus() == "completed" {
				if run.GetConclusion() == "success" {
					ecrPushSuccessful = true
				}
			} else {
				// Workflow still in progress — fetch the jobs
				jobs, _, err := ghClient.Actions.ListWorkflowJobs(ctx, config.Owner, config.Repo, run.GetID(), &github.ListWorkflowJobsOptions{
					ListOptions: github.ListOptions{PerPage: 100},
				})
				if err != nil {
					return constants.AutoDeployStatusInProgress, nil, fmt.Errorf("failed to fetch workflow jobs: %w", err)
				}

				var jobStatuses []WorkflowJobStatus
				for _, job := range jobs.Jobs {
					jobStatuses = append(jobStatuses, WorkflowJobStatus{
						Name:       job.GetName(),
						Status:     job.GetStatus(),
						Conclusion: job.GetConclusion(),
					})
				}

				return constants.AutoDeployStatusInProgress, jobStatuses, nil
			}
			break
		}
	}

	if !ecrPushSuccessful {
		return constants.AutoDeployStatusFailed, nil, fmt.Errorf("ECR Push workflow failed. Aborting deployments")
	}

	return constants.AutoDeployStatusCompleted, nil, nil
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
