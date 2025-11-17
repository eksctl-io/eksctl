package delete

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	capabilityactions "github.com/weaveworks/eksctl/pkg/actions/capability"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func deleteCapabilityCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"capability",
		"Delete capabilities",
		"",
	)

	capability := &api.Capability{}
	configureDeleteCapabilityCmd(cmd, capability)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewDeleteCapabilityLoader(cmd, capability).Load(); err != nil {
			return err
		}
		return doDeleteCapability(cmd, capability)
	}
}

func doDeleteCapability(cmd *cmdutils.Cmd, capability *api.Capability) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.ProviderConfig.WaitTimeout)
	defer cancel()

	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	// Use capabilities from config file if available, otherwise use single capability
	var capabilities []capabilityactions.CapabilitySummary
	if len(cmd.ClusterConfig.Capabilities) > 0 {
		for _, cap := range cmd.ClusterConfig.Capabilities {
			capabilities = append(capabilities, capabilityactions.CapabilitySummary{
				Capability: cap,
			})
		}

	} else {
		capabilities = []capabilityactions.CapabilitySummary{{
			Capability: *capability,
		}}
	}

	stackManager := clusterProvider.NewStackManager(cmd.ClusterConfig)

	capabilityRemover := capabilityactions.NewRemover(cmd.ClusterConfig.Metadata.Name, stackManager, clusterProvider.AWSProvider.EKS(), cmd.ProviderConfig.WaitTimeout)
	// Delete capabilities and IAM role
	return capabilityRemover.Delete(ctx, capabilities)

}

func configureDeleteCapabilityCmd(cmd *cmdutils.Cmd, capability *api.Capability) {
	cmd.FlagSetGroup.InFlagSet("Capability", func(fs *pflag.FlagSet) {
		fs.StringVar(&capability.Name, "name", "", "Name of the capability to delete")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}
