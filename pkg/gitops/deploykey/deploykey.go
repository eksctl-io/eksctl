package deploykey

import (
	"context"
	"os"
	"strings"

	"github.com/fluxcd/go-git-providers/github"
	"github.com/fluxcd/go-git-providers/gitprovider"
	"github.com/pkg/errors"
	giturls "github.com/whilp/git-urls"
)

const githubTokenVariable = "GITHUB_TOKEN"

// GetDeployKeyClient provides a DeployKeyClient for the repo at the given URL
func GetDeployKeyClient(ctx context.Context, url string) (gitprovider.DeployKeyClient, error) {
	repo, err := giturls.Parse(url)
	if err != nil {
		return nil, err
	}
	if strings.Contains(repo.Host, "github.com") {
		githubToken := os.Getenv(githubTokenVariable)
		gh, err := github.NewClient(github.WithOAuth2Token(githubToken))
		if err != nil {
			return nil, err
		}
		ownerRepo := strings.Split(repo.Path, "/")
		if len(ownerRepo) != 2 {
			return nil, errors.Errorf("couldn't handle URL as github.com URL: %s", url)
		}
		rep, err := gh.UserRepositories().Get(ctx, gitprovider.UserRepositoryRef{
			UserRef: gitprovider.UserRef{
				Domain:    "github.com",
				UserLogin: ownerRepo[0],
			},
			RepositoryName: ownerRepo[1],
		})
		if err != nil {
			return nil, err
		}
		return rep.DeployKeys(), nil
	}
	return nil, nil
}
