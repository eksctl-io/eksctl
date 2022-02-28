package utils

import (
	"context"
	"time"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

const enableKMSTimeout = (1 * time.Hour) + (20 * time.Minute)

func enableSecretsEncryptionWithHandler(cmd *cmdutils.Cmd, handler func(*cmdutils.Cmd, bool) error) {
	cfg := api.NewClusterConfig()
	cfg.SecretsEncryption = &api.SecretsEncryption{}
	cmd.ClusterConfig = cfg

	cmd.SetDescription("enable-secrets-encryption", "Enable secrets encryption", "Enable secrets encryption on a cluster")

	var encryptExistingSecrets bool

	cmd.CobraCommand.RunE = func(_ *cobra.Command, args []string) error {
		cmd.NameArg = cmdutils.GetNameArg(args)
		if err := cmdutils.NewUtilsKMSLoader(cmd).Load(); err != nil {
			return err
		}
		return handler(cmd, encryptExistingSecrets)
	}

	cmdutils.AddCommonFlagsForAWS(cmd.FlagSetGroup, &cmd.ProviderConfig, false)

	cmd.FlagSetGroup.InFlagSet("General", func(fs *pflag.FlagSet) {
		cmdutils.AddClusterFlag(fs, cfg.Metadata)
		cmdutils.AddRegionFlag(fs, &cmd.ProviderConfig)
		cmdutils.AddConfigFileFlag(fs, &cmd.ClusterConfigFile)
		cmdutils.AddApproveFlag(fs, cmd)
		cmdutils.AddTimeoutFlagWithValue(fs, &cmd.ProviderConfig.WaitTimeout, enableKMSTimeout)
		fs.StringVar(&cmd.ClusterConfig.SecretsEncryption.KeyARN, "key-arn", "", "KMS key ARN")
		fs.BoolVar(&encryptExistingSecrets, "encrypt-existing-secrets", true, "Encrypt all existing secrets with the new KMS key")
	})

}

func enableSecretsEncryptionCmd(cmd *cmdutils.Cmd) {
	enableSecretsEncryptionWithHandler(cmd, doEnableSecretsEncryption)
}

func doEnableSecretsEncryption(cmd *cmdutils.Cmd, encryptExistingSecrets bool) error {
	clusterConfig := cmd.ClusterConfig

	ctl, err := cmd.NewProviderForExistingCluster()
	if err != nil {
		return err
	}

	if cmd.ClusterConfig.Metadata.Name == "" {
		return cmdutils.ErrMustBeSet(cmdutils.ClusterNameFlag(cmd))
	}

	if err := ctl.RefreshClusterStatus(clusterConfig); err != nil {
		return err
	}

	if err := api.ValidateSecretsEncryption(clusterConfig); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), ctl.Provider.WaitTimeout())
	defer cancel()

	if err := ctl.EnableKMSEncryption(ctx, clusterConfig); err != nil {
		return err
	}

	if encryptExistingSecrets {
		logger.Info("updating all Secret resources to apply KMS encryption")
		clientSet, err := ctl.NewStdClientSet(clusterConfig)
		if err != nil {
			return err
		}
		if err := kubernetes.RefreshSecrets(ctx, clientSet.CoreV1()); err != nil {
			return errors.Wrap(err, "error updating secrets")
		}
		logger.Info("KMS encryption applied to all Secret resources")
	}

	return nil
}
