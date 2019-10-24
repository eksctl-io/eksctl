package enable

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
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
	opts := ConfigureRepositoryCmd(cmd)
	cmd.SetRunFuncWithNameArg(func() error {
		return Repository(cmd, opts)
	})
}

// ConfigureRepositoryCmd configures the provided command object so that it can
// process CLI options and ClusterConfig file, to prepare for the "enablement"
// of the configured repository & cluster.
func ConfigureRepositoryCmd(cmd *cmdutils.Cmd) *flux.InstallOpts {
	var opts flux.InstallOpts
	cmd.FlagSetGroup.InFlagSet("Enable repository", func(fs *pflag.FlagSet) {
		cmdutils.AddCommonFlagsForGit(fs, &opts.GitOptions)

		fs.StringSliceVar(&opts.GitPaths, gitPaths, []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.GitLabel, gitLabel, "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.GitFluxPath, gitFluxPath, "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.StringVar(&opts.Namespace, namespace, "flux",
			"Cluster namespace where to install Flux, the Helm Operator and Tiller")
		fs.BoolVar(&opts.WithHelm, withHelm, true,
			"Install the Helm Operator and Tiller")
		fs.BoolVar(&opts.Amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVar(&cmd.ClusterConfig.Metadata.Name, "cluster", "", "name of the EKS cluster to enable gitops on")
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.Timeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	cmd.ProviderConfig.WaitTimeout = opts.Timeout
	return &opts
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
	if opts.GitOptions.PrivateSSHKeyPath != "" && cfg.Git.Repo.PrivateSSHKeyPath != "" {
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

// Repository enables GitOps on the configured repository.
func Repository(cmd *cmdutils.Cmd, opts *flux.InstallOpts) error {
	if err := cmdutils.NewInstallFluxLoader(cmd).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig

	if cfg.HasGitOpsOptions() {
		if err := validateGitOpsOptions(cfg, opts); err != nil {
			return err
		}
		optsFromCfg, err := flux.NewInstallOptsFrom(cfg.Git, opts.Timeout)
		if err != nil {
			return err
		}
		opts.CopyFrom(optsFromCfg)
	} else {
		cmdutils.ValidateGitOptions(&opts.GitOptions)
	}

	k8sClientSet, k8sRestConfig, err := KubernetesClientAndConfigFrom(cmd)
	if err != nil {
		return err
	}

	installer := flux.NewInstaller(k8sRestConfig, k8sClientSet, opts)
	userInstructions, err := installer.Run(context.Background())
	logger.Info(userInstructions)
	return err
}
