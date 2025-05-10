package constants

const (
	CommonFlagDescNamespace       = "Target Kubernetes namespace"
	CommonFlagDescDeploymentNames = "Comma-separated list of deployment names"
	CommonFlagDescPodNames        = "List of pod names (required)"
	CommonFlagDescTag             = "Image tag to deploy (required)"
	CommonFlagDescLocation        = "Location to the file (default: current directory)"

	CommonFlagDescWatchFrequency = "Frequency in seconds to update %s information"
	CommonFlagDescPodLevel       = "Display information at pod level"

	AutoDeployGitTagToWatch         = "Git tag to watch"
	AutoDeployGithubRepository      = "GitHub repository name"
	AutoDeployGithubRepositoryOwner = "GitHub repository owner (e.g., 'octocat' for https://github.com/octocat/repo)"
)
