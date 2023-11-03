package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/actions/accessentry"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func deleteAccessEntryCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()

	cmd.SetDescription(
		"accessentry",
		"Delete access entry(ies)",
		"",
		"accessentries",
	)

	var principalARN string
	cmd.FlagSetGroup.InFlagSet("AccessEntry", func(fs *pflag.FlagSet) {
		fs.StringVar(&principalARN, "principal-arn", "", "principal ARN to which the access entry is associated")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doDeleteAccessEntry(cmd, principalARN)
	}
}

func doDeleteAccessEntry(cmd *cmdutils.Cmd, principalARN string) error {
	cmd.ClusterConfig.AccessConfig.AccessEntries = []api.AccessEntry{{}}
	if principalARN != "" {
		cmd.ClusterConfig.AccessConfig.AccessEntries[0].PrincipalARN = api.MustParseARN(principalARN)
	}

	if err := cmdutils.NewDeleteAccessEntryLoader(cmd).Load(); err != nil {
		return err
	}

	ctx := context.Background()
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	accessEntryManager := accessentry.NewRemover(
		cmd.ClusterConfig,
		clusterProvider.NewStackManager(cmd.ClusterConfig),
		clusterProvider.AWSProvider.EKS(),
	)

	if err = accessEntryManager.Delete(ctx, cmd.ClusterConfig.AccessConfig.AccessEntries); err != nil {
		return err
	}

	return nil
}
