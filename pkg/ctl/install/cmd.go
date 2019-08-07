package install

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

// Command will create the `install` commands
func Command(flagGrouping *cmdutils.FlagGrouping) *cobra.Command {
	verbCmd := cmdutils.NewVerbCmd("install", "Install components in a cluster", "")

	cmdutils.AddResourceCmd(flagGrouping, verbCmd, installFluxCmd)

	return verbCmd
}

type installFluxOpts struct {
	gitURL              string
	gitBranch           string
	gitPaths            []string
	gitLabel            string
	gitUser             string
	gitEmail            string
	gitFluxPath         string
	namespace           string
	tillerNamespace     string
	timeout             time.Duration
	amend               bool
	noHelmOp            bool
	noTiller            bool
	tillerHost          string
	helmOpTLSCertFile   string
	helmOpTLSKeyFile    string
	helmOpTLSCACertFile string
}

func installFluxCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"flux",
		"Bootstrap Flux, installing it in the cluster and initializing its manifests in the specified Git repository",
		"",
	)
	var opts installFluxOpts
	cmd.SetRunFuncWithNameArg(func() error {
		installer, err := newFluxInstaller(context.Background(), cmd, &opts)
		if err != nil {
			return err
		}
		return installer.run(context.Background())
	})
	cmd.FlagSetGroup.InFlagSet("Flux installation", func(fs *pflag.FlagSet) {
		fs.StringVar(&opts.gitURL, "git-url", "",
			"URL of the Git repository to be used by Flux, e.g. git@github.com:<your username>/flux-get-started")
		fs.StringVar(&opts.gitBranch, "git-branch", "master",
			"Git branch to be used by Flux")
		fs.StringSliceVar(&opts.gitPaths, "git-paths", []string{},
			"Relative paths within the Git repo for Flux to locate Kubernetes manifests")
		fs.StringVar(&opts.gitLabel, "git-label", "flux",
			"Git label to keep track of Flux's sync progress; overrides both --git-sync-tag and --git-notes-ref")
		fs.StringVar(&opts.gitUser, "git-user", "Flux",
			"Username to use as Git committer")
		fs.StringVar(&opts.gitEmail, "git-email", "",
			"Email to use as Git committer")
		fs.StringVar(&opts.gitFluxPath, "git-flux-subdir", "flux/",
			"Directory within the Git repository where to commit the Flux manifests")
		fs.StringVar(&opts.namespace, "namespace", "flux",
			"Cluster namespace where to install Flux and the Helm Operator")
		fs.StringVar(&opts.tillerNamespace, "tiller-namespace", "kube-system",
			"Cluster namespace where to install Tiller")
		fs.BoolVar(&opts.noHelmOp, "no-helmop", false,
			"Do not install the Flux Helm Operator (implies --no-tiller)")
		fs.BoolVar(&opts.noTiller, "no-tiller", false,
			"Do not install Tiller. If the Helm Operator is installed, make sure to appropriately set --tiller-host and --helmop-tls-* options")
		fs.StringVar(&opts.tillerHost, "tiller-host", "",
			"Only applicable when --no-tiller is set. Hostname to use in order to connect to Tiller")
		fs.StringVar(&opts.helmOpTLSKeyFile, "helmop-tls-key-path", "",
			"Only applicable when --no-tiller is set. Path to TLS client key to connect to Tiller. Empty means no TLS.")
		fs.StringVar(&opts.helmOpTLSCertFile, "helmop-tls-cert-path", "",
			"Only applicable when --no-tiller is set. Path to TLS client certificate to connect to Tiller. Empty means no TLS.")
		fs.StringVar(&opts.helmOpTLSCACertFile, "helmop-tls-ca-cert-path", "",
			"Only applicable when --no-tiller is set. Path to CA certificate to verify Tiller's certificate. Empty means no TLS.")
		fs.BoolVar(&opts.amend, "amend", false,
			"Stop to manually tweak the Flux manifests before pushing them to the Git repository")
	})
	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddNameFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlagWithValue(fs, &opts.timeout, 20*time.Second)
	})
	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, cmd.ProviderConfig, false)
	cmd.ProviderConfig.WaitTimeout = opts.timeout
}
