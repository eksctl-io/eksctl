package deploykey

import (
	"context"
	"os"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
)

const githubTokenVariable = "GITHUB_TOKEN"

// GetDeployKeyClient provides a DeployKeyClient for the repo at the given URL
func GetDeployKeyClient(ctx context.Context, url string) (gitprovider.DeployKeyClient, error) {
	ref, err := gitprovider.ParseUserRepositoryURL(url)
	if err != nil {
		return nil, err
	}
	switch ref.Domain {
	case "github.com":
		githubToken := os.Getenv(githubTokenVariable)
		gh, err := github.NewClient(github.WithOAuth2Token(githubToken))
		if err != nil {
			return nil, err
		}
		rep, err := gh.UserRepositories().Get(ctx, *ref)
		if err != nil {
			return nil, err
		}
		return rep.DeployKeys(), nil
	}
	return nil, nil
}
