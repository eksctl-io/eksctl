package enable

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/git"
	"github.com/weaveworks/eksctl/pkg/gitops"
	"github.com/weaveworks/eksctl/pkg/gitops/fileprocessor"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
)

type options struct {
	gitOptions      git.Options
	profileNameArg  string
	profileRevision string
}

func (opts options) validate() error {
	if opts.profileNameArg == "" {
		return errors.New("please supply a valid Quick Start profile URL or name")
	}
	if err := opts.gitOptions.ValidateURL(); err != nil {
		return errors.Wrap(err, "please supply a valid --git-url argument")
	}
	if err := opts.gitOptions.ValidateEmail(); err != nil {
		return errors.Wrap(err, "please supply a valid --git-email argument")
	}
	if err := opts.gitOptions.ValidatePrivateSSHKeyPath(); err != nil {
		return errors.Wrap(err, "please supply a valid --git-private-ssh-key-path argument")
	}
	return nil
}

func enableProfileCmd(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("profile", "Set up Flux and deploy the components from the selected Quick Start profile.", "")

	var opts options

	cmd.SetRunFuncWithNameArg(func() error {
		return doEnableProfile(cmd, opts)
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		fs.StringVarP(&opts.profileNameArg, "name", "", "", "name or URL of the Quick Start profile. For example, app-dev.")
		fs.StringVarP(&opts.profileRevision, "revision", "", "master", "revision of the Quick Start profile.")
		fs.StringVarP(&opts.gitOptions.URL, "git-url", "", "", "SSH URL of the Git repository that will contain the cluster components, e.g. git@github.com:<github_org>/<repo_name>")
		fs.StringVarP(&opts.gitOptions.Branch, "git-branch", "", "master", "Git branch")
		fs.StringVar(&opts.gitOptions.User, "git-user", "Flux", "Username to use as Git committer")
		fs.StringVar(&opts.gitOptions.Email, "git-email", "", "Email to use as Git committer")
		fs.StringVar(&opts.gitOptions.PrivateSSHKeyPath, "git-private-ssh-key-path", "",
			"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
		fs.StringVar(&cfg.Metadata.Name, "cluster", "", "name of the EKS cluster to add the Quick Start profile to")

		requiredFlags := []string{"git-url", "git-email"}
		for _, f := range requiredFlags {
			if err := cobra.MarkFlagRequired(fs, f); err != nil {
				logger.Critical("unexpected error: %v", err)
				os.Exit(1)
			}
		}

		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, 20*time.Second)
	})

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
}

func doEnableProfile(cmd *cmdutils.Cmd, opts options) error {
	if cmd.NameArg != "" && opts.profileNameArg != "" {
		return cmdutils.ErrClusterFlagAndArg(cmd, cmd.NameArg, opts.profileNameArg)
	}
	if cmd.NameArg != "" {
		opts.profileNameArg = cmd.NameArg
	}
	if err := opts.validate(); err != nil {
		return err
	}

	profileRepoURL, err := repoURLForQuickstart(opts.profileNameArg)
	if err != nil {
		return errors.Wrap(err, "please supply a valid Quick Start profile name or URL")
	}

	if err := cmdutils.NewEnableProfileLoader(cmd).Load(); err != nil {
		return err
	}
	cfg := cmd.ClusterConfig
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

	// Create the flux installer. It will clone the user's repository in a temporary directory.
	fluxOpts := flux.InstallOpts{
		GitOptions:  opts.gitOptions,
		Namespace:   "flux",
		GitFluxPath: "flux/",
		WithHelm:    true,
		Timeout:     cmd.ProviderConfig.WaitTimeout,
	}
	fluxInstaller := flux.NewInstaller(k8sRestConfig, k8sClientSet, &fluxOpts)

	processor := &fileprocessor.GoTemplateProcessor{
		Params: fileprocessor.NewTemplateParameters(cmd.ClusterConfig),
	}

	// Create the profile generator. It will output the processed templates into a new "/base" directory into the user's repo
	usersRepoName, err := git.RepoName(opts.gitOptions.URL)
	if err != nil {
		return err
	}
	dir, err := ioutil.TempDir("", usersRepoName)
	logger.Debug("Directory %s will be used to clone the configuration repository and install the profile", dir)
	usersRepoDir := filepath.Join(dir, usersRepoName)
	profileOutputPath := filepath.Join(usersRepoDir, "base")

	profile := &gitops.Profile{
		Processor: processor,
		Path:      profileOutputPath,
		GitOpts: git.Options{
			URL:    profileRepoURL,
			Branch: opts.profileRevision,
		},
		GitCloner: git.NewGitClient(git.ClientParams{}),
		FS:        afero.NewOsFs(),
		IO:        afero.Afero{Fs: afero.NewOsFs()},
	}

	// A git client that operates in the user's repo
	gitClient := git.NewGitClient(git.ClientParams{
		PrivateSSHKeyPath: opts.gitOptions.PrivateSSHKeyPath,
	})

	gitOps := gitops.Applier{
		UserRepoPath:     usersRepoDir,
		UsersRepoOpts:    opts.gitOptions,
		GitClient:        gitClient,
		ProfileGenerator: profile,
		FluxInstaller:    fluxInstaller,
		ClusterConfig:    cmd.ClusterConfig,
		QuickstartName:   opts.profileNameArg,
	}

	if err = gitOps.Run(context.Background()); err != nil {
		return err
	}
	os.RemoveAll(dir) // Only clean up if the command completely successfully, for more convenient debugging.
	return nil
}

func repoURLForQuickstart(quickstartArgument string) (string, error) {
	if git.IsGitURL(quickstartArgument) {
		return quickstartArgument, nil
	}
	if quickstartArgument == "app-dev" {
		return "https://github.com/weaveworks/eks-quickstart-app-dev", nil
	}
	if quickstartArgument == "appmesh" {
		return "https://github.com/weaveworks/eks-appmesh-profile", nil
	}
	return "", fmt.Errorf("invalid URL or unknown Quick Start profile %s ", quickstartArgument)
}
