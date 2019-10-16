package enable

import (
	"context"
	"fmt"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
	"github.com/weaveworks/eksctl/pkg/utils/file"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	gitURL               = "git-url"
	gitBranch            = "git-branch"
	gitUser              = "git-user"
	gitEmail             = "git-email"
	gitPaths             = "git-paths"
	gitFluxPath          = "git-flux-subdir"
	gitLabel             = "git-label"
	gitPrivateSSHKeyPath = "git-private-ssh-key-path"
	namespace            = "namespace"
	withHelm             = "with-helm"
)

func enableRepo(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"repo",
		"Set up a repo for gitops, installing Flux in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	var opts flux.InstallOpts
	cmd.SetRunFuncWithNameArg(func() error {
		if err := cmdutils.NewInstallFluxLoader(cmd).Load(); err != nil {
			return err
		}
		cfg := cmd.ClusterConfig

		if cfg.HasGitOpsOptions() {
			if err := validateGitOpsOptions(cfg, &opts); err != nil {
				return err
			}
			optsFromCfg, err := flux.NewInstallOptsFrom(cfg.Git, opts.Timeout)
			if err != nil {
				return err
			}
			opts.CopyFrom(optsFromCfg)
		} else {
			validateInstallOpts(&opts)
		}

		ctl, err := cmd.NewCtl()
		if err != nil {
			return err
		}

		if err := ctl.CheckAuth(); err != nil {
			return err
		}
		if ok, err := ctl.CanOperate(cfg); !ok {
			return err
		}
		kubernetesClientConfigs, err := ctl.NewClient(cfg)
		if err != nil {
			return err
		}
		k8sConfig := kubernetesClientConfigs.Config

		k8sRestConfig, err := clientcmd.NewDefaultClientConfig(*k8sConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
		if err != nil {
			return errors.Wrap(err, "cannot create Kubernetes client configuration")
		}
		k8sClientSet, err := kubeclient.NewForConfig(k8sRestConfig)
		if err != nil {
			return errors.Errorf("cannot create Kubernetes client set: %s", err)
		}

		installer := flux.NewInstaller(k8sRestConfig, k8sClientSet, &opts)
		userInstructions, err := installer.Run(context.Background())
		logger.Info(userInstructions)
		return err
	})

	cmd.FlagSetGroup.InFlagSet("Flux installation", func(fs *pflag.FlagSet) {
		fs.StringVar(&opts.GitOptions.URL, gitURL, "",
			"SSH URL of the Git repository to be used by Flux, e.g. git@github.com:<github_org>/<repo_name>")
		fs.StringVar(&opts.GitOptions.Branch, gitBranch, "master",
			"Git branch to be used by Flux")
		fs.StringSliceVar(&opts.GitPaths, gitPaths, []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.GitLabel, gitLabel, "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.GitOptions.User, gitUser, "Flux",
			"Username to use as Git committer")
		fs.StringVar(&opts.GitOptions.Email, gitEmail, "",
			"Email to use as Git committer")
		fs.StringVar(&opts.GitFluxPath, gitFluxPath, "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.StringVar(&opts.GitPrivateSSHKeyPath, gitPrivateSSHKeyPath, "",
			"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
		fs.StringVar(&opts.Namespace, namespace, "flux",
			"Cluster namespace where to install Flux, the Helm Operator and Tiller")
		fs.BoolVar(&opts.WithHelm, withHelm, true,
			"Install the Helm Operator and Tiller")
		fs.BoolVar(&opts.Amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "EKS cluster name")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.Timeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	cmd.ProviderConfig.WaitTimeout = opts.Timeout
}

func validateGitOpsOptions(cfg *api.ClusterConfig, opts *flux.InstallOpts) error {
	if opts.GitOptions.URL != "" && cfg.Git.Repo.URL != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitURL)
	}
	if opts.GitOptions.Branch != "" && cfg.Git.Repo.Branch != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitBranch)
	}
	if opts.GitOptions.User != "" && cfg.Git.Repo.User != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitUser)
	}
	if opts.GitOptions.Email != "" && cfg.Git.Repo.Email != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitEmail)
	}
	if len(opts.GitPaths) > 0 && len(cfg.Git.Repo.Paths) > 0 {
		return cmdutils.ErrCannotUseWithConfigFile(gitPaths)
	}
	if opts.GitFluxPath != "" && cfg.Git.Repo.FluxPath != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitFluxPath)
	}
	if opts.GitLabel != "" && cfg.Git.Operator.Label != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitLabel)
	}
	if opts.GitPrivateSSHKeyPath != "" && cfg.Git.Repo.PrivateSSHKeyPath != "" {
		return cmdutils.ErrCannotUseWithConfigFile(gitPrivateSSHKeyPath)
	}
	if opts.Namespace != "" && cfg.Git.Operator.Namespace != "" {
		return cmdutils.ErrCannotUseWithConfigFile(namespace)
	}
	if opts.WithHelm && !cfg.Git.Operator.WithHelm {
		return cmdutils.ErrCannotUseWithConfigFile(withHelm)
	}
	if err := api.ValidateGit(cfg.Git); err != nil {
		return err
	}
	return nil
}

func validateInstallOpts(opts *flux.InstallOpts) error {
	if err := opts.GitOptions.ValidateURL(); err != nil {
		return errors.Wrapf(err, "please supply a valid --%s argument", gitURL)
	}
	if opts.GitOptions.Email == "" {
		return fmt.Errorf("please supply a valid --%s argument", gitEmail)
	}
	if opts.GitPrivateSSHKeyPath != "" && !file.Exists(opts.GitPrivateSSHKeyPath) {
		return fmt.Errorf("please supply a valid --%s argument", gitPrivateSSHKeyPath)
	}
	return nil
}
