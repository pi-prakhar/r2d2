package ghservice

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

var sharedGitHubClient *github.Client

// GetSharedGitHubClient returns a singleton instance of the GitHub client.
// It initializes the client if it hasn't been created yet.
func GetSharedGitHubClient(accessToken string) *github.Client {
	if sharedGitHubClient == nil {
		ctx := context.Background()
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		tokenClient := oauth2.NewClient(ctx, tokenSource)
		sharedGitHubClient = github.NewClient(tokenClient)
	}
	return sharedGitHubClient
}
