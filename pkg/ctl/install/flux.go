package install

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
	"github.com/weaveworks/eksctl/pkg/utils/file"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func installFluxCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"flux",
		"Bootstrap Flux, installing it in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	var opts flux.InstallOpts
	cmd.SetRunFuncWithNameArg(func() error {
		if err := opts.GitOptions.ValidateURL(); err != nil {
			return errors.Wrapf(err, "please supply a valid --git-url argument")
		}
		if opts.GitOptions.Email == "" {
			return errors.New("please supply a valid --git-email argument")
		}
		if opts.GitPrivateSSHKeyPath != "" && !file.Exists(opts.GitPrivateSSHKeyPath) {
			return errors.New("please supply a valid --git-private-ssh-key-path argument")
		}

		if err := cmdutils.NewMetadataLoader(cmd).Load(); err != nil {
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
		if err := ctl.RefreshClusterConfig(cfg); err != nil {
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

		installer := flux.NewInstaller(context.Background(), k8sRestConfig, k8sClientSet, &opts)
		return installer.Run(context.Background())
	})

	cmd.FlagSetGroup.InFlagSet("Flux installation", func(fs *pflag.FlagSet) {
		fs.StringVar(&opts.GitOptions.URL, "git-url", "",
			"SSH URL of the Git repository to be used by Flux, e.g. git@github.com:<github_org>/<repo_name>")
		fs.StringVar(&opts.GitOptions.Branch, "git-branch", "master",
			"Git branch to be used by Flux")
		fs.StringSliceVar(&opts.GitPaths, "git-paths", []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.GitLabel, "git-label", "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.GitOptions.User, "git-user", "Flux",
			"Username to use as Git committer")
		fs.StringVar(&opts.GitOptions.Email, "git-email", "",
			"Email to use as Git committer")
		fs.StringVar(&opts.GitFluxPath, "git-flux-subdir", "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.StringVar(&opts.GitPrivateSSHKeyPath, "git-private-ssh-key-path", "",
			"Optional path to the private SSH key to use with Git, e.g. ~/.ssh/id_rsa")
		fs.StringVar(&opts.Namespace, "namespace", "flux",
			"Cluster namespace where to install Flux, the Helm Operator and Tiller")
		fs.BoolVar(&opts.WithHelm, "with-helm", true,
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
