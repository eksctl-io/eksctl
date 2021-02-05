// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
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
	if !strings.Contains(repo.Host, github.DefaultDomain) {
		return nil, errors.New("only GitHub URLs are currently supported")
	}
	githubToken := os.Getenv(githubTokenVariable)
	if githubToken == "" {
		return nil, errors.Errorf("%s not set", githubTokenVariable)
	}
	gh, err := github.NewClient(github.WithOAuth2Token(githubToken))
	if err != nil {
		return nil, err
	}
	ownerRepo := strings.Split(repo.Path, "/")
	if len(ownerRepo) != 2 {
		return nil, errors.New("couldn't understand path of URL")
	}
	repoName := strings.TrimSuffix(ownerRepo[1], ".git")
	rep, err := gh.UserRepositories().Get(ctx, gitprovider.UserRepositoryRef{
		UserRef: gitprovider.UserRef{
			Domain:    github.DefaultDomain,
			UserLogin: ownerRepo[0],
		},
		RepositoryName: repoName,
	})
	if err != nil {
		return nil, err
	}
	return rep.DeployKeys(), nil
}
