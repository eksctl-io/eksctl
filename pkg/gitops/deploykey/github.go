package deploykey

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v31/github"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"golang.org/x/oauth2"
)

const (
	EnvVarGitHubToken = "GITHUB_TOKEN"
)

type GitHubProvider struct {
	cluster     *api.ClusterMeta
	owner, repo string
	readOnly    bool
	githubToken string
}

func (p *GitHubProvider) Put(ctx context.Context, fluxSSHKey PublicKey) error {
	gh := p.getGitHubAPIClient(ctx)

	logger.Info("creating GitHub deploy key from Flux SSH public key")

	title := p.getDeployKeyTitle()

	key, _, err := gh.Repositories.CreateKey(ctx, p.owner, p.repo, &github.Key{
		Key:      &fluxSSHKey.Key,
		Title:    &title,
		ReadOnly: &p.readOnly,
	})

	if err != nil {
		return err
	}

	logger.Info("%s configured with Flux SSH public key\n%s", *key.Title, fluxSSHKey.Key)

	return nil
}

func (p *GitHubProvider) Delete(ctx context.Context) error {
	gh := p.getGitHubAPIClient(ctx)

	logger.Info("deleting GitHub deploy key")

	title := p.getDeployKeyTitle()

	keys, _, err := gh.Repositories.ListKeys(ctx, p.owner, p.repo, &github.ListOptions{})
	if err != nil {
		return err
	}

	var keyID int64

	for _, key := range keys {
		if key.GetTitle() == title {
			keyID = key.GetID()

			break
		}
	}

	if keyID == 0 {
		logger.Info("skipped deleting GitHub deploy key %q: The key does not exist. Probably you've already deleted it?")

		return nil
	}

	if _, err := gh.Repositories.DeleteKey(ctx, p.owner, p.repo, keyID); err != nil {
		return err
	}

	logger.Info("deleted GitHub deploy key %s", title)

	return nil
}

func (p *GitHubProvider) getGitHubAPIClient(ctx context.Context) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	gh := github.NewClient(tc)

	return gh
}

func getGitHubOwnerRepoFromRepoURL(repoURL string) (string, string, bool) {
	if repoURL == "" {
		return "", "", false
	}

	sshFull := regexp.MustCompile(`ssh://git@github.com/([^/]+)/([^.]+).git`)
	sshShort := regexp.MustCompile(`git@github.com:([^/]+)/([^.]+).git`)

	patterns := []*regexp.Regexp{
		sshFull,
		sshShort,
	}

	for _, p := range patterns {
		m := p.FindStringSubmatch(repoURL)
		if len(m) == 3 {
			return m[1], m[2], true
		}
	}

	return "", "", false
}

func (p *GitHubProvider) getDeployKeyTitle() string {
	return fmt.Sprintf("eksctl-flux-%s-%s", p.cluster.Region, p.cluster.Name)
}
