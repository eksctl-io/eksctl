package deploykey

import (
	"context"
	"os"

	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type GitProvider interface {
	Put(ctx context.Context, fluxSSHKey PublicKey) error
	Delete(ctx context.Context) error
}

func ForCluster(cluster *v1alpha5.ClusterConfig) GitProvider {
	var (
		repoURL  string
		readOnly bool
	)

	if git := cluster.Git; git != nil {
		if repo := git.Repo; repo != nil {
			repoURL = repo.URL
		}

		readOnly = git.Operator.ReadOnly
	}

	if repoURL == "" {
		return nil
	}

	if owner, repo, ok := getGitHubOwnerRepoFromRepoURL(repoURL); !ok {
		logger.Info("skipped managing GitHub deploy key for URL %s: Only `git@github.com:OWNER/REPO.git` is accepted for automatic deploy key creation", repoURL)
	} else if githubToken := os.Getenv(EnvVarGitHubToken); githubToken == "" {
		logger.Info("GITHUB_TOKEN is not set. Please set it so that eksctl is able to create and delete GitHub deploy key from Flux SSH public key")
	} else {
		return &GitHubProvider{
			cluster:     cluster.Metadata,
			githubToken: githubToken,
			readOnly:    readOnly,
			owner:       owner,
			repo:        repo,
		}
	}

	return nil
}

func Put(ctx context.Context, cluster *v1alpha5.ClusterConfig, fluxSSHKey PublicKey) (bool, error) {
	p := ForCluster(cluster)

	if p == nil {
		return false, nil
	}

	return true, p.Put(ctx, fluxSSHKey)
}

func Delete(ctx context.Context, cluster *v1alpha5.ClusterConfig) error {
	p := ForCluster(cluster)

	if p == nil {
		return nil
	}

	return p.Delete(ctx)
}

// PublicKey represents a public SSH key as it is returned by flux
type PublicKey struct {
	Key string
}
