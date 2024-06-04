package utils

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func describeAddonConfigurationCmd(cmd *cmdutils.Cmd) {
	cmd.SetDescription(
		"describe-addon-configuration",
		"Output the configuration JSON schema for an addon",
		"",
	)

	var (
		addonName    string
		addonVersion string
	)
	cmd.FlagSetGroup.InFlagSet("Addon", func(fs *pflag.FlagSet) {
		fs.StringVar(&addonName, "name", "", "Addon name")
		fs.StringVar(&addonVersion, "version", "", "Addon version")
		_ = cobra.MarkFlagRequired(fs, "name")
		_ = cobra.MarkFlagRequired(fs, "version")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddTimeoutFlag(fs, &cmd.ProviderConfig.WaitTimeout)
	})
	cmdutils.AddCommonFlagsForAWS(cmd, &cmd.ProviderConfig, false)

	cmd.CobraCommand.RunE = func(_ *cobra.Command, _ []string) error {
		return describeAddonConfiguration(cmd, addonName, addonVersion)
	}
}

func describeAddonConfiguration(cmd *cmdutils.Cmd, addonName, addonVersion string) error {
	ctx := context.Background()
	clusterProvider, err := eks.New(ctx, &cmd.ProviderConfig, nil)
	if err != nil {
		return err
	}

	addonConfig, err := clusterProvider.AWSProvider.EKS().DescribeAddonConfiguration(ctx, &awseks.DescribeAddonConfigurationInput{
		AddonName:    aws.String(addonName),
		AddonVersion: aws.String(addonVersion),
	})
	if err != nil {
		return fmt.Errorf("error describing addon configuration: %w", err)
	}
	if addonConfig.ConfigurationSchema == nil {
		return fmt.Errorf("no configuration schema found for %s@%s", addonName, addonVersion)
	}

	var schema interface{}
	if err := json.Unmarshal([]byte(*addonConfig.ConfigurationSchema), &schema); err != nil {
		return fmt.Errorf("unmarshalling retrieved addon configuration schema: %w", err)
	}
	config, err := json.MarshalIndent(struct {
		Schema      any                                   `json:"configurationSchema"`
		PodIDConfig []types.AddonPodIdentityConfiguration `json:"podIdentityConfiguration"`
	}{
		Schema:      schema,
		PodIDConfig: addonConfig.PodIdentityConfiguration,
	}, "", "\t")
	if err != nil {
		return fmt.Errorf("marshalling retrieved addon configuration: %w", err)
	}
	fmt.Println(string(config))

	return nil
}
