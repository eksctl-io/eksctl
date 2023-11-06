package utils

import (
	"context"
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
)

func updateAuthenticationMode(cmd *cmdutils.Cmd) {
	cfg := api.NewClusterConfig()
	cmd.ClusterConfig = cfg

	cmd.SetDescription("update-authentication-mode", "Updates the authentication mode for a cluster", "")

	var authenticationMode string
	cmd.FlagSetGroup.InFlagSet("Authentication mode", func(fs *pflag.FlagSet) {
		fs.StringVar(&authenticationMode, "authentication-mode", "", "authentication mode of the cluster")
	})

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cmd.ClusterConfig.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
	})

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		return doUpdateAuthenticationMode(cmd, authenticationMode)
	}
}

func doUpdateAuthenticationMode(cmd *cmdutils.Cmd, authenticationMode string) error {
	cmd.ClusterConfig.AccessConfig.AuthenticationMode = ekstypes.AuthenticationMode(authenticationMode)
	if err := cmdutils.NewUtilsUpdateAuthenticationModeLoader(cmd).Load(); err != nil {
		return err
	}

	ctx := context.Background()
	clusterProvider, err := cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return err
	}

	logger.Info("setting cluster's authentication mode to %s", authenticationMode)
	clusterName := cmd.ClusterConfig.Metadata.Name
	output, err := clusterProvider.AWSProvider.EKS().UpdateClusterConfig(ctx, &awseks.UpdateClusterConfigInput{
		Name: aws.String(clusterName),
		AccessConfig: &ekstypes.UpdateAccessConfigRequest{
			AuthenticationMode: cmd.ClusterConfig.AccessConfig.AuthenticationMode,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update cluster config: %v", err)
	}

	updateWaiter := waiter.NewUpdateWaiter(clusterProvider.AWSProvider.EKS(), func(options *waiter.UpdateWaiterOptions) {
		options.RetryAttemptLogMessage = fmt.Sprintf("waiting for update %q in cluster %q to complete", *output.Update.Id, clusterName)
	})
	err = updateWaiter.Wait(ctx, &awseks.DescribeUpdateInput{
		Name:     &clusterName,
		UpdateId: output.Update.Id,
	}, clusterProvider.AWSProvider.WaitTimeout())

	switch e := err.(type) {
	case *waiter.UpdateFailedError:
		if e.Status == string(ekstypes.UpdateStatusCancelled) {
			return fmt.Errorf("request to update cluster authentication mode was cancelled: %s", e.UpdateError)
		}
		return fmt.Errorf("failed to update cluster authentication mode: %s", e.UpdateError)

	case nil:
		logger.Info("authentication mode was successfully updated to %s on cluster %s", cmd.ClusterConfig.AccessConfig.AuthenticationMode, clusterName)
		return nil

	default:
		return err
	}
}
