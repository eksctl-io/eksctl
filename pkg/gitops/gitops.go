package gitops

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
)

// Setup sets up gitops in a repository for a cluster. Installs flux, helm and a quickstart profile into the cluster
func Setup(k8sRestConfig *rest.Config, k8sClientSet kubeclient.Interface, cfg *api.ClusterConfig, timeout time.Duration) error {
	installer, err := flux.NewInstaller(k8sRestConfig, k8sClientSet, cfg, timeout)
	if err != nil {
		return err
	}

	fluxIsInstalled, err := installer.IsFluxInstalled()
	if err != nil {
		// Continue with installation
		logger.Warning(err.Error())
	} else if fluxIsInstalled {
		logger.Warning("found existing flux deployment in namespace %q. Skipping installation",
			cfg.Git.Operator.Namespace)
		return nil
	}

	userInstructions, err := installer.Run(context.Background())
	if err != nil {
		return errors.Wrapf(err, "unable to install flux")
	}

	err = InstallProfile(cfg)
	if err != nil {
		return err
	}

	logger.Info(userInstructions)
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

	err := profileGen.Install(cfg)
	if err != nil {
		return errors.Wrapf(err, "unable to install bootstrap profile")
	}

	return nil
}
