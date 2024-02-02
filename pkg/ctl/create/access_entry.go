package create

import (
	"context"

	"github.com/kris-nova/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/accessentry"
	accessentryactions "github.com/weaveworks/eksctl/pkg/actions/accessentry"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils/filter"
)

func createAccessEntryCmdWithRunFunc(cmd *cmdutils.Cmd, runFunc func(cmd *cmdutils.Cmd) error) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"accessentry",
		"Create access entries",
		"",
	)

	accessEntry := &api.AccessEntry{}
	configureCreateAccessEntryCmd(cmd, accessEntry)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewCreateAccessEntryLoader(cmd, accessEntry).Load(); err != nil {
			return err
		}
		return runFunc(cmd)
	}
}

func createAccessEntryCmd(cmd *cmdutils.Cmd) {
	createAccessEntryCmdWithRunFunc(cmd, doCreateAccessEntry)
}

func doCreateAccessEntry(cmd *cmdutils.Cmd) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.ProviderConfig.WaitTimeout)
	defer cancel()

	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}
	accessEntry := &accessentry.Service{
		ClusterStateGetter: clusterProvider,
	}
	if !accessEntry.IsEnabled() {
		return accessentry.ErrDisabledAccessEntryAPI
	}
	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)
	accessEntryFilter := &filter.AccessEntry{
		Lister:      stackManager,
		ClusterName: cmd.ClusterConfig.Metadata.Name,
	}
	accessEntries, err := accessEntryFilter.FilterOutExistingStacks(ctx, cmd.ClusterConfig.AccessConfig.AccessEntries)
	if err != nil {
		return err
	}
	if len(accessEntries) == 0 {
		logger.Info("access entries already up-to-date")
		return nil
	}
	accessEntryCreator := &accessentryactions.Creator{
		ClusterName:  cmd.ClusterConfig.Metadata.Name,
		StackCreator: stackManager,
	}
	return accessEntryCreator.Create(ctx, accessEntries)
}

func configureCreateAccessEntryCmd(cmd *cmdutils.Cmd, accessEntry *api.AccessEntry) {
	cmd.FlagSetGroup.InFlagSet("Access Entry", func(fs *pflag.FlagSet) {
		fs.VarP(&accessEntry.PrincipalARN, "principal-arn", "", "Principal ARN")
		fs.StringVar(&accessEntry.Type, "type", "", "Type of Access Entry")
		fs.StringSliceVar(&accessEntry.KubernetesGroups, "kubernetes-groups", nil, "A set of Kubernetes groups to map to the principal ARN")
		fs.StringVar(&accessEntry.KubernetesUsername, "kubernetes-username", "", "A Kubernetes username to map to the principal ARN")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}
