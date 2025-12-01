package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	capabilityactions "github.com/weaveworks/eksctl/pkg/actions/capability"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
)

func updateCapabilityCmd(cmd *cmdutils.Cmd) {
	cmd.ClusterConfig = api.NewClusterConfig()
	cmd.SetDescription(
		"capability",
		"Update capabilities",
		"",
	)

	capability := &api.Capability{}
	configureUpdateCapabilityCmd(cmd, capability)
}

func doUpdateCapability(cmd *cmdutils.Cmd, capability *api.Capability, attachPolicyStr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), cmd.ProviderConfig.WaitTimeout)
	defer cancel()

	var capabilities []api.Capability

	// Use capabilities from config file if available, otherwise use single capability from flags
	if len(cmd.ClusterConfig.Capabilities) > 0 {
		capabilities = cmd.ClusterConfig.Capabilities
	} else {
		// Parse attach policy if provided
		if attachPolicyStr != "" {
			if err := parseAttachPolicy(attachPolicyStr, capability); err != nil {
				return err
			}
		}
		capabilities = []api.Capability{*capability}
	}

	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	capabilityUpdater := capabilityactions.NewUpdater(cmd.ClusterConfig.Metadata.Name, clusterProvider.AWSProvider.EKS())

	for _, cap := range capabilities {
		logger.Info("updating capability %s", cap.Name)
	}
	return capabilityUpdater.Update(ctx, capabilities)
}

func configureUpdateCapabilityCmd(cmd *cmdutils.Cmd, capability *api.Capability) {
	var attachPolicyStr string
	cmd.FlagSetGroup.InFlagSet("Capability", func(fs *pflag.FlagSet) {
		fs.StringVar(&capability.Name, "name", "", "Name of the capability")
		fs.StringVar(&capability.Type, "type", "", "Type of the capability (ACK, KRO, ARGOCD)")
		fs.StringVar(&capability.RoleARN, "role-arn", "", "IAM role ARN for the capability (optional if IAM policies are provided)")
		fs.StringToStringVar(&capability.Tags, "tags", nil, "Tags to apply to the capability")
		fs.StringSliceVar(&capability.AttachPolicyARNs, "attach-policy-arns", nil, "List of IAM policy ARNs to attach to the role")
		fs.StringVar(&attachPolicyStr, "attach-policy", "", "Inline IAM policy document to attach to the role (JSON string)")
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewCreateCapabilityLoader(cmd, capability).Load(); err != nil {
			return err
		}
		return doUpdateCapability(cmd, capability, attachPolicyStr)
	}

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})

	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)
}

func parseAttachPolicy(policyStr string, capability *api.Capability) error {
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(policyStr), &policy); err != nil {
		return fmt.Errorf("invalid JSON in attach-policy: %w", err)
	}
	capability.AttachPolicy = policy
	return nil
}
