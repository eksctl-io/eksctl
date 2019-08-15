package install

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	kubeclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/gitops/flux"
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
		if opts.GitURL == "" {
			return errors.New("please supply a valid --git-url argument")
		}
		if opts.GitEmail == "" {
			return errors.New("please supply a valid --git-email argument")
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
		fs.StringVar(&opts.GitURL, "git-url", "",
			"URL of the Git repository to be used by Flux, e.g. git@github.com:<github_org>/flux-get-started")
		fs.StringVar(&opts.GitBranch, "git-branch", "master",
			"Git branch to be used by Flux")
		fs.StringSliceVar(&opts.GitPaths, "git-paths", []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.GitLabel, "git-label", "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.GitUser, "git-user", "Flux",
			"Username to use as Git committer")
		fs.StringVar(&opts.GitEmail, "git-email", "",
			"Email to use as Git committer")
		fs.StringVar(&opts.GitFluxPath, "git-flux-subdir", "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.StringVar(&opts.Namespace, "namespace", "flux",
			"Cluster namespace where to install Flux, the Helm Operator and Tiller")
		fs.BoolVar(&opts.WithHelm, "with-helm", true,
			"Install the Helm Operator and Tiller")
		fs.BoolVar(&opts.Amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.Timeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	cmd.ProviderConfig.WaitTimeout = opts.Timeout
}
