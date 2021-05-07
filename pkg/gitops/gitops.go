package gitops

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/weaveworks/eksctl/pkg/actions/flux"
	"github.com/weaveworks/eksctl/pkg/actions/repo"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/deploykey"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
)

// DefaultPodReadyTimeout is the time it will wait for Flux and Helm Operator to become ready
const DefaultPodReadyTimeout = 5 * time.Minute

type FluxInstaller interface {
	Run() error
}

// Setup sets up gitops in a repository for a cluster.
func Setup(kubeconfigPath string, k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, cfg *api.ClusterConfig, timeout time.Duration) error {
	installer, profilesSupported, err := newFluxInstaller(kubeconfigPath, k8sRestConfig, k8sClientSet, cfg, timeout)
	if err != nil {
		return errors.Wrapf(err, "could not initialise Flux installer")
	}

	if err := installer.Run(); err != nil {
		return err
	}

	if profilesSupported {
		return InstallProfile(cfg)
	}

	return nil
}

// InstallProfile installs the bootstrap profile in the user's repo if it's specified in the cluster config
func InstallProfile(cfg *api.ClusterConfig) error {
	if !cfg.HasBootstrapProfile() {
		logger.Debug("no bootstrap profiles configure. Skipping...")
		return nil
	}

	gitCfg := cfg.Git

	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: gitCfg.Repo.PrivateSSHKeyPath,
	})

	profileGen := &Profile{
		Processor: &fileprocessor.GoTemplateProcessor{
			Params: fileprocessor.NewTemplateParameters(cfg),
		},
		UserRepoGitClient: gitClient,
		ProfileCloner: git.NewGitClient(git.ClientParams{
			PrivateSSHKeyPath: gitCfg.Repo.PrivateSSHKeyPath,
		}),
		FS: afero.NewOsFs(),
		IO: afero.Afero{Fs: afero.NewOsFs()},
	}

	if err := profileGen.Install(cfg); err != nil {
		return errors.Wrapf(err, "unable to install bootstrap profile")
	}

	return nil
}

// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
// DeleteKey deletes the authorized SSH key for the gitops repo if gitops are configured
// Will not fail if the key was not previously authorized
func DeleteKey(cfg *api.ClusterConfig) error {
	if !cfg.HasGitopsRepoConfigured() {
		return nil
	}

	ctx := context.Background()
	deployKeyClient, err := deploykey.GetDeployKeyClient(ctx, cfg.Git.Repo.URL)
	if err != nil {
		logger.Warning(
			"could not find git provider implementation for url %q: %q. Skipping deletion of authorized SSH key",
			cfg.Git.Repo.URL,
			err.Error(),
		)
		return nil
	}

	clusterKeyTitle := repo.KeyTitle(*cfg.Metadata)
	logger.Info("deleting SSH key %q from repo %q", clusterKeyTitle, cfg.Git.Repo.URL)

	key, err := deployKeyClient.Get(ctx, clusterKeyTitle)
	if err != nil {
		return errors.Wrapf(err, "unable to find SSH key")
	}
	if err := key.Delete(ctx); err != nil {
		return errors.Wrapf(err, "unable to delete authorized key")
	}
	return nil
}

func newFluxInstaller(kubeconfigPath string, k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, cfg *api.ClusterConfig, timeout time.Duration) (FluxInstaller, bool, error) {
	var (
		installer         FluxInstaller
		profilesSupported bool
		err               error
	)

	if cfg.GitOps != nil && cfg.GitOps.Flux != nil {
		installer, err = flux.New(k8sClientSet, cfg.GitOps)
		profilesSupported = false
		logger.Info("gitops configuration detected, setting installer to Flux v2")
	} else {
		// FLUX V1 DEPRECATION NOTICE. https://github.com/weaveworks/eksctl/issues/2963
		installer, err = repo.New(k8sRestConfig, k8sClientSet, cfg, timeout)
		profilesSupported = true
		logger.Info("git.repo configuration detected, setting installer to Flux v1")
	}

	return installer, profilesSupported, err
}
